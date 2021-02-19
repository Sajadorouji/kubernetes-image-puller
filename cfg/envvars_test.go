package cfg

import (
	corev1 "k8s.io/api/core/v1"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestProcessImagesEnvVar(t *testing.T) {
	type testcase struct {
		name   string
		images string
		want   map[string]string
	}

	testcases := []testcase{
		{
			name:   "one image",
			images: "che-theia=quay.io/eclipse/che-theia:nightly",
			want:   map[string]string{"che-theia": "quay.io/eclipse/che-theia:nightly"},
		},
		{
			name:   "three images",
			images: "image1=my/image1:dev;image2=my/image2:next;image3=my/image3:stage",
			want: map[string]string{
				"image1": "my/image1:dev",
				"image2": "my/image2:next",
				"image3": "my/image3:stage",
			},
		},
	}

	for _, c := range testcases {
		t.Run(c.name, func(t *testing.T) {
			defer os.Clearenv()
			os.Setenv("IMAGES", c.images)
			got := processImagesEnvVar()

			if d := cmp.Diff(c.want, got); d != "" {
				t.Errorf("(-want, +got): %s", d)
			}
		})
	}
}

func TestProcessNodeSElectorEnvVar(t *testing.T) {
	type testcase struct {
		name              string
		nodeSelector      string
		isNodeSelectorSet bool
		want              map[string]string
	}

	testcases := []testcase{
		{
			name:              "default node selector, NODE_SELECTOR set",
			nodeSelector:      "{}",
			isNodeSelectorSet: true,
			want:              map[string]string{},
		},
		{
			name:              "compute type, NODE_SELECTOR set",
			nodeSelector:      "{\"type\": \"compute\"}",
			isNodeSelectorSet: true,
			want: map[string]string{
				"type": "compute",
			},
		},
		{
			name:              "default env var, NODE_SELECTOR not set",
			nodeSelector:      "{\"this\": \"shouldn't be set\"}",
			isNodeSelectorSet: false,
			want:              map[string]string{},
		},
	}

	for _, c := range testcases {
		t.Run(c.name, func(t *testing.T) {
			defer os.Clearenv()
			if c.isNodeSelectorSet {
				os.Setenv("NODE_SELECTOR", c.nodeSelector)
			}
			got := processNodeSelectorEnvVar()

			if d := cmp.Diff(c.want, got); d != "" {
				t.Errorf("(-want, +got): %s", d)
			}
		})
	}
}
func TestProcessNodeTolerationEnvVar(t *testing.T) {
	type testcase struct {
		name                string
		nodeToleration      string
		isNodeTolerationSet bool
		want                []corev1.Toleration
	}

	testcases := []testcase{
		{
			name:              "default node selector, NODE_TOLERATION set",
			nodeToleration:      "[]",
			isNodeTolerationSet: true,
			want:                []corev1.Toleration{},
		},
		{
			name:              "compute type, NODE_TOLERATION set",
			nodeToleration:      "[{\"effect\": \"NoSchedule\",\"key\": \"test\",\"operator\": \"Exist\",\"value\": \"yes\"}]",
			isNodeTolerationSet: true,
			want: []corev1.Toleration{
				{
					Key:      "test",
					Operator: "Exist",
					Value:    "yes",
					Effect:   "NoSchedule",
				},
			},
		},
		{
			name:                "default env var, NODE_TOLERATION not set",
			nodeToleration:      "[{\"this\": \"shouldn't be set\"}]",
			isNodeTolerationSet: false,
			want:                []corev1.Toleration{},
		},
	}

	for _, c := range testcases {
		t.Run(c.name, func(t *testing.T) {
			defer os.Clearenv()
			if c.isNodeTolerationSet {
				os.Setenv("NODE_TOLERATION", c.nodeToleration)
			}
			got := processNodeTolerationEnvVar()

			if d := cmp.Diff(c.want, got); d != "" {
				t.Errorf("(-want, +got): %s", d)
			}
		})
	}
}


func TestGetEnvVarOrDefaultBool(t *testing.T) {
	defer os.Clearenv()
	os.Setenv("DEFINED_ENV_VAR", "true")
	assert.True(t, getEnvVarOrDefaultBool("DEFINED_ENV_VAR", false), "When a variable is defined it should return its value")
	assert.True(t, getEnvVarOrDefaultBool("UNDEFINED_ENV_VAR", true), "When a variable is not defined it should return the default value")
}

func TestGetEnvVarOrDefault(t *testing.T) {
	defer os.Clearenv()
	os.Setenv("DEFINED_ENV_VAR", "foo")
	assert.Equal(t, getEnvVarOrDefault("DEFINED_ENV_VAR", "bar"), "foo", "When a variable is defined it should return it's set value")
	assert.Equal(t, getEnvVarOrDefault("UNDEFINED_ENV_VAR", "bar"), "bar", "When a variable is undefined it should return the default value")
}
