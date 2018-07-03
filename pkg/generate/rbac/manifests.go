/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rbac

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManifestOptions represent options for generating the RBAC manifests.
type ManifestOptions struct {
	InputDir  string
	OutputDir string
	Name      string
	Labels    map[string]string
}

// Validate validates the input options.
func (o *ManifestOptions) Validate() error {
	if _, err := os.Stat(o.InputDir); err != nil {
		return fmt.Errorf("invalid input directory '%s' %v", o.InputDir, err)
	}
	if _, err := os.Stat(o.OutputDir); err != nil {
		return fmt.Errorf("invalid output directory '%s' %v", o.OutputDir, err)
	}
	return nil
}

// Generate generates RBAC manifests by parsing the RBAC annotations in Go source
// files specified in the input directory.
func Generate(o *ManifestOptions) error {
	if err := o.Validate(); err != nil {
		return err
	}

	rules, err := ParseDir(o.InputDir)
	if err != nil {
		return fmt.Errorf("failed to parse the input dir %v", err)
	}
	if len(rules) == 0 {
		return nil
	}
	roleManifest, err := getClusterRoleManifest(rules, o)
	if err != nil {
		return fmt.Errorf("failed to generate role manifest %v", err)
	}

	roleBindingManifest, err := getClusterRoleBindingManifest(o)
	if err != nil {
		return fmt.Errorf("failed to generate role binding manifests %v", err)
	}

	roleManifestFile := filepath.Join(o.OutputDir, "rbac_role.yaml")
	if err := ioutil.WriteFile(roleManifestFile, roleManifest, 0666); err != nil {
		return fmt.Errorf("failed to write role manifest YAML file %v", err)
	}

	roleBindingManifestFile := filepath.Join(o.OutputDir, "rbac_role_binding.yaml")
	if err := ioutil.WriteFile(roleBindingManifestFile, roleBindingManifest, 0666); err != nil {
		return fmt.Errorf("failed to write role manifest YAML file %v", err)
	}
	return nil
}

func getClusterRoleManifest(rules []rbacv1.PolicyRule, o *ManifestOptions) ([]byte, error) {
	role := rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   o.Name + "-role",
			Labels: o.Labels,
		},
		Rules: rules,
	}
	return yaml.Marshal(role)
}

func getClusterRoleBindingManifest(o *ManifestOptions) ([]byte, error) {
	rolebinding := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-rolebinding", o.Name),
			Labels: o.Labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Name:      "default",
				Namespace: fmt.Sprintf("%v-system", o.Name),
				Kind:      "ServiceAccount",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Name:     fmt.Sprintf("%v-role", o.Name),
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	return yaml.Marshal(rolebinding)
}
