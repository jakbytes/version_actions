package semver

import "testing"

func TestPrereleaseIdentifier(t *testing.T) {
	v, _ := NewVersion("1.0.0-alpha.1")
	expected := "alpha"
	if id := v.PrereleaseIdentifier(); id != expected {
		t.Errorf("Expected %s, got %s", expected, id)
	}
}

func TestPrereleaseVersionNumber(t *testing.T) {
	v, _ := NewVersion("1.0.0-alpha.1")
	expected := 1
	num, _ := v.PrereleaseVersionNumber()
	if num != expected {
		t.Errorf("Expected %d, got %d", expected, num)
	}
}

func TestIsPrerelease(t *testing.T) {
	v, _ := NewVersion("1.0.0-alpha.1")
	if !v.IsPrerelease() {
		t.Errorf("Expected true, got false")
	}

	v2, _ := NewVersion("1.0.0")
	if v2.IsPrerelease() {
		t.Errorf("Expected false, got true")
	}
}

func TestNewVersion(t *testing.T) {
	_, err := NewVersion("invalid")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	_, err2 := NewVersion("1.0.0")
	if err2 != nil {
		t.Errorf("Expected nil, got error")
	}
}

func TestIncMajor(t *testing.T) {
	v, _ := NewVersion("1.2.3")
	v2 := v.IncMajor()
	if v2.String() != "2.0.0" {
		t.Errorf("Expected 2.0.0, got %s", v2.String())
	}
}

func TestMustParse(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParse did not panic on invalid version")
		}
	}()
	MustParse("invalid-version")
}

func TestIncMinor(t *testing.T) {
	v, _ := NewVersion("1.2.3")
	v2 := v.IncMinor()
	expected := "1.3.0"
	if v2.String() != expected {
		t.Errorf("Expected %s, got %s", expected, v2.String())
	}
}

func TestIncPatch(t *testing.T) {
	v, _ := NewVersion("1.2.3")
	v2 := v.IncPatch()
	expected := "1.2.4"
	if v2.String() != expected {
		t.Errorf("Expected %s, got %s", expected, v2.String())
	}
}

func TestAsPrereleaseVersion(t *testing.T) {
	v, _ := NewVersion("1.2.3")
	v2, err := v.AsPrereleaseVersion("beta", 2)
	if err != nil {
		t.Errorf("Error creating prerelease version: %v", err)
	}
	expected := "1.2.3-beta.2"
	if v2.String() != expected {
		t.Errorf("Expected %s, got %s", expected, v2.String())
	}
}

func TestGreaterThan(t *testing.T) {
	v1, _ := NewVersion("1.2.3")
	v2, _ := NewVersion("1.2.4")
	if !v2.GreaterThan(v1) {
		t.Errorf("Expected %s to be greater than %s", v2.String(), v1.String())
	}

	if v1.GreaterThan(v2) {
		t.Errorf("Expected %s not to be greater than %s", v1.String(), v2.String())
	}
}

func TestLessThan(t *testing.T) {
	v1, _ := NewVersion("1.2.3")
	v2, _ := NewVersion("1.2.4")
	if !v1.LessThan(v2) {
		t.Errorf("Expected %s to be less than %s", v1.String(), v2.String())
	}

	if v2.LessThan(v1) {
		t.Errorf("Expected %s not to be less than %s", v2.String(), v1.String())
	}
}
