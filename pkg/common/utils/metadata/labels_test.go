// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package metadata

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_LabelsTest(t *testing.T) {
	labels := Labels(map[string]string{"test": "test", "test1": "test1"})
	ls := NewLabels(labels)
	addLabels := Labels(map[string]string{"testadd": "testadd"})
	ls.AddLabel(addLabels)
	ls.Add("addkey", "addkey")
	if len(ls) != 4 {
		t.Errorf("test labels have not right number.")
	}
}

func Test_Annotation(t *testing.T) {
	anno := Annotations(map[string]string{"test": "test", "test1": "test1"})
	nanno := NewAnnotations(anno)
	addNno := Annotations(map[string]string{"addtest": "addtest"})
	nanno.AddAnnotation(addNno)
	nanno.Add("addtestkey", "addtestkey")
	if len(nanno) != 4 {
		t.Error("test annotation have not right number.")
	}
}

func Test_MergeMetadata(t *testing.T) {
	nom := &metav1.ObjectMeta{
		ResourceVersion: "2",
		Finalizers:      []string{"newFinalizer"},
		Labels:          map[string]string{"newtest": "newtest", "newtest1": "newtest1"},
		Annotations:     map[string]string{"newannotest": "newannotest", "newannotest1": "newannotest1"},
		OwnerReferences: []metav1.OwnerReference{{
			APIVersion: "v1",
			Kind:       "Statefulset",
			Name:       "test",
		}},
	}

	oom := metav1.ObjectMeta{
		ResourceVersion: "1",
		Finalizers:      []string{"oldFinalizer"},
		Labels:          map[string]string{"oldtest": "oldtest", "oldtest1": "oldtest1"},
		Annotations:     map[string]string{"oldannotest": "oldannotest", "oldannotest1": "oldannotest1"},
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "v1",
				Kind:       "Statefulset",
				Name:       "test",
			}, {
				APIVersion: "v1",
				Kind:       "DorisCluster",
				Name:       "testcluster",
			},
		},
	}

	MergeMetadata(nom, oom)
	if len(nom.Finalizers) != 2 {
		t.Errorf("mergeMetadata merge finalizer not right.")
	}
	if len(nom.Labels) != 4 {
		t.Errorf("mergeMetadata merge labels not right.")
	}
	if len(nom.OwnerReferences) != 1 {
		t.Errorf("mergeMetadata merge ownerReferences not right.")
	}
	if len(nom.Annotations) != 4 {
		t.Errorf("meregeMetadata merge annotations not right.")
	}
}
