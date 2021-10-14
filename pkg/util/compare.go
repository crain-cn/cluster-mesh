package util

import corev1 "k8s.io/api/core/v1"

func IsEqualResourceRequirement(require, require2 corev1.ResourceRequirements) bool {
	return require.Limits.Cpu().Equal(*require2.Limits.Cpu()) && require.Limits.Memory().Equal(*require2.Limits.Memory()) && require.Requests.Memory().Equal(*require2.Requests.Memory()) && require.Requests.Cpu().Equal(*require2.Requests.Cpu())
}
