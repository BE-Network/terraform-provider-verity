package lifecycle

import (
	"regexp"
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func secretSwitchpoint(t *testing.T, base, lines string) string {
	t.Helper()
	base = regexp.MustCompile(`(?m)^  ssh_key_or_password = .*\n`).ReplaceAllString(base, "")
	return strings.Replace(base, "  name = \"diffsp\"\n", "  name = \"diffsp\"\n"+lines, 1)
}

func secretMock(t *testing.T) (*mock.MockServer, ResourceCoverageEntry) {
	t.Helper()
	entry := coverageEntry(t, "verity_switchpoint")
	ms := mock.NewMockServer(entry.Mode)
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatal(err)
	}
	ms.SetSecretMasking(entry.APIPath, map[string]string{
		"ssh_key_or_password": "write_only",
	})
	return ms, entry
}

func TestSecretFieldsRoundTripWithTheV66API(t *testing.T) {
	ms, entry := secretMock(t)
	defer ms.Close()
	base := switchpointBase(t)
	provider := mock.ProviderConfig(ms.URL(), entry.Mode)

	plaintext := func(password string) string {
		return provider + secretSwitchpoint(t, base,
			"  ssh_key_or_password = \""+password+"\"\n  ssh_key_or_password_encrypted_auto_assigned_ = false\n")
	}
	encrypted := provider + secretSwitchpoint(t, base,
		"  ssh_key_or_password = \"NewPass1\"\n  ssh_key_or_password_encrypted = \"$enc$abc\"\n  ssh_key_or_password_encrypted_auto_assigned_ = false\n")

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), plaintext("test-password-1")) },
				Config:    plaintext("test-password-1"),
				Check:     fwresource.TestCheckResourceAttr("verity_switchpoint.test", "ssh_key_or_password", "test-password-1"),
			},
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), plaintext("NewPass1")) },
				Config:    plaintext("NewPass1"),
				Check:     fwresource.TestCheckResourceAttr("verity_switchpoint.test", "ssh_key_or_password", "NewPass1"),
			},
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), encrypted) },
				Config:    encrypted,
				Check:     fwresource.TestCheckResourceAttr("verity_switchpoint.test", "ssh_key_or_password_encrypted", "$enc$abc"),
			},
		},
	})

	var sawNewPassword bool
	for _, patch := range ms.GetRequestsByMethodAndPath("PATCH", entry.APIPath) {
		body := canonical(t, patch.Body)
		if strings.Contains(strings.ToLower(body), mock.MaskedSecretPlaceholder) {
			t.Errorf("the placeholder was sent: %s", body)
		}
		if strings.Contains(body, `"ssh_key_or_password":"NewPass1"`) {
			sawNewPassword = true
		}
	}
	if !sawNewPassword {
		t.Error("the changed password was not sent in a PATCH")
	}
	for _, put := range ms.GetRequestsByMethodAndPath("PUT", entry.APIPath) {
		if strings.Contains(strings.ToLower(canonical(t, put.Body)), mock.MaskedSecretPlaceholder) {
			t.Errorf("the placeholder was sent on create: %s", canonical(t, put.Body))
		}
	}
}

func TestEncryptedPasswordFromTheAPIIsKept(t *testing.T) {
	ms, entry := secretMock(t)
	defer ms.Close()
	ms.SetPostPutEnrichment(entry.APIPath, "switchpoint", "diffsp", map[string]interface{}{"ssh_key_or_password_encrypted": "$enc$from-api"})
	config := mock.ProviderConfig(ms.URL(), entry.Mode) + secretSwitchpoint(t, switchpointBase(t),
		"  ssh_key_or_password = \"test-password-1\"\n  ssh_key_or_password_encrypted_auto_assigned_ = false\n")

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:    config,
				Check: fwresource.ComposeTestCheckFunc(
					fwresource.TestCheckResourceAttr("verity_switchpoint.test", "ssh_key_or_password", "test-password-1"),
					fwresource.TestCheckResourceAttr("verity_switchpoint.test", "ssh_key_or_password_encrypted", "$enc$from-api"),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestFlagOffKeepsThePasswordWithoutSendingIt(t *testing.T) {
	ms, entry := secretMock(t)
	defer ms.Close()
	base := switchpointBase(t)
	provider := mock.ProviderConfig(ms.URL(), entry.Mode)
	config := func(flag string) string {
		return provider + secretSwitchpoint(t, base, "  ssh_key_or_password_encrypted_auto_assigned_ = "+flag+"\n")
	}

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config("true")) }, Config: config("true")},
			{PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config("false")) }, Config: config("false")},
		},
	})

	var flagPatch string
	for _, patch := range ms.GetRequestsByMethodAndPath("PATCH", entry.APIPath) {
		if body := canonical(t, patch.Body); strings.Contains(body, `"ssh_key_or_password_encrypted_auto_assigned_":false`) {
			flagPatch = body
		}
	}
	if flagPatch == "" {
		t.Fatal("no PATCH turned the flag off")
	}
	if strings.Contains(flagPatch, `"ssh_key_or_password"`) || strings.Contains(flagPatch, `"ssh_key_or_password_encrypted"`) {
		t.Fatalf("turning the flag off on 6.6 needs no password, and none must be sent: %s", flagPatch)
	}
}
