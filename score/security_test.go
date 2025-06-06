package score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestContainerSeccompMissing(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-seccomp-no-annotation.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeWarning)

}

func TestPodSecurityContextPresent(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-seccomp-securecontext-ok.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeAllOK)
}

func TestContainerSecurityContextSeccompPresent(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-seccomp-container-securecontext-ok.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeAllOK)
}

func TestPodSecurityContextSeccompAbsent(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-seccomp-securecontext-warning.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeWarning)
}

func TestContainerSecurityContextSeccompAbsent(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-seccomp-container-securecontext-warning.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeWarning)
}

func TestContainerSeccompMissingNotRunByDefault(t *testing.T) {
	t.Parallel()
	skipped := wasSkipped(t, []ks.NamedReader{testFile("pod-seccomp-no-annotation.yaml")}, nil, nil, "Container Seccomp Profile")
	assert.True(t, skipped)
}

func TestContainerSeccompAllGood(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-seccomp-profile"] = struct{}{}

	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-seccomp-annotated.yaml")}, nil, &config.RunConfiguration{

		EnabledOptionalTests: structMap,
	}, "Container Seccomp Profile", scorecard.GradeAllOK)
}

func TestContainerSeccompAllGoodAnnotation(t *testing.T) {
	t.Parallel()

	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-seccomp-annotated-annotation-optional.yaml")}, nil, &config.RunConfiguration{
		UseOptionalChecksAnnotation: true,
	}, "Container Seccomp Profile", scorecard.GradeAllOK)
}

func TestContainerSecurityContextUserGroupIDAllGood(t *testing.T) {
	t.Parallel()
	structMap := make(map[string]struct{})
	structMap["container-security-context-user-group-id"] = struct{}{}
	c := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-all-good.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Security Context User Group ID", scorecard.GradeAllOK)
	assert.Empty(t, c)
}

func TestContainerSecurityContextUserGroupIDLowGroup(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-user-group-id"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-low-group-id.yaml")}, nil, &config.RunConfiguration{

		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context User Group ID", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The container running with a low group ID",
		Description: "A groupid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsGroup to a value > 10000",
	})
}

func TestContainerSecurityContextUserGroupIDLowUser(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-user-group-id"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-low-user-id.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context User Group ID", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The container is running with a low user ID",
		Description: "A userid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsUser to a value > 10000",
	})
}

func TestContainerSecurityContextUserGroupIDNoSecurityContext(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-user-group-id"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-nosecuritycontext.yaml")}, nil, &config.RunConfiguration{

		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context User Group ID", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "Container has no configured security context",
		Description: "Set securityContext to run the container in a more secure context.",
	})
}

func TestContainerSecurityContextPrivilegedAllGood(t *testing.T) {
	t.Parallel()
	structMap := make(map[string]struct{})
	structMap["container-security-context-privileged"] = struct{}{}
	c := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-all-good.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Security Context Privileged", scorecard.GradeAllOK)
	assert.Empty(t, c)
}

func TestContainerSecurityContextPrivilegedPrivileged(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-privileged"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-privileged.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context Privileged", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The container is privileged",
		Description: "Set securityContext.privileged to false. Privileged containers can access all devices on the host, and grants almost the same access as non-containerized processes on the host.",
	})
}

func TestContainerSecurityContextReadOnlyRootFilesystemAllGood(t *testing.T) {
	t.Parallel()
	structMap := make(map[string]struct{})
	structMap["container-security-context-readonlyrootfilesystem"] = struct{}{}
	c := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-all-good.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: structMap,
	}, "Container Security Context ReadOnlyRootFilesystem", scorecard.GradeAllOK)
	assert.Empty(t, c)
}

func TestContainerSecurityContextReadOnlyRootFilesystemWriteable(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-readonlyrootfilesystem"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-writeablerootfilesystem.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context ReadOnlyRootFilesystem", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "The pod has a container with a writable root filesystem",
		Description: "Set securityContext.readOnlyRootFilesystem to true",
	})
}

func TestContainerSecurityContextReadOnlyRootFilesystemNoSecurityContext(t *testing.T) {
	t.Parallel()
	optionalChecks := make(map[string]struct{})
	optionalChecks["container-security-context-readonlyrootfilesystem"] = struct{}{}
	comments := testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("pod-security-context-nosecuritycontext.yaml")}, nil, &config.RunConfiguration{
		EnabledOptionalTests: optionalChecks,
	}, "Container Security Context ReadOnlyRootFilesystem", scorecard.GradeCritical)
	assert.Contains(t, comments, scorecard.TestScoreComment{
		Path:        "foobar",
		Summary:     "Container has no configured security context",
		Description: "Set securityContext to run the container in a more secure context.",
	})
}
