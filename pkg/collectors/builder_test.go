// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package collectors

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"
	koptions "k8s.io/kube-state-metrics/pkg/options"
	"k8s.io/kube-state-metrics/pkg/whiteblacklist"
)

var (
	ctx context.Context = context.TODO()
)

func TestBuilder_WithApiserver(t *testing.T) {
	type fields struct {
		apiserver         string
		kubeconfig        string
		namespaces        koptions.NamespaceList
		ctx               context.Context
		enabledCollectors []string
		whiteBlackList    whiteBlackLister
	}
	type args struct {
		apiserver string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "apiserver",
			fields: fields{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
			},
			args: args{
				apiserver: "apiserver",
			},
			want: &Builder{
				apiserver:         "apiserver",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				apiserver:         tt.fields.apiserver,
				kubeconfig:        tt.fields.kubeconfig,
				namespaces:        tt.fields.namespaces,
				ctx:               tt.fields.ctx,
				enabledCollectors: tt.fields.enabledCollectors,
				whiteBlackList:    tt.fields.whiteBlackList,
			}
			if got := b.WithApiserver(tt.args.apiserver); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.WithApiserver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithKubeConfig(t *testing.T) {
	type fields struct {
		apiserver         string
		kubeconfig        string
		namespaces        koptions.NamespaceList
		ctx               context.Context
		enabledCollectors []string
		whiteBlackList    whiteBlackLister
	}
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "kubeconfig",
			fields: fields{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
			},
			args: args{
				kubeconfig: "kubeconfig",
			},
			want: &Builder{
				apiserver:         "",
				kubeconfig:        "kubeconfig",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				apiserver:         tt.fields.apiserver,
				kubeconfig:        tt.fields.kubeconfig,
				namespaces:        tt.fields.namespaces,
				ctx:               tt.fields.ctx,
				enabledCollectors: tt.fields.enabledCollectors,
				whiteBlackList:    tt.fields.whiteBlackList,
			}
			if got := b.WithKubeConfig(tt.args.kubeconfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.WithKubeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithEnabledCollectors(t *testing.T) {
	type fields struct {
		apiserver         string
		kubeconfig        string
		namespaces        koptions.NamespaceList
		ctx               context.Context
		enabledCollectors []string
		whiteBlackList    whiteBlackLister
	}
	type args struct {
		c []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "collectors",
			fields: fields{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{},
			},
			args: args{
				c: []string{"col1", "col2"},
			},
			want: &Builder{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				apiserver:         tt.fields.apiserver,
				kubeconfig:        tt.fields.kubeconfig,
				namespaces:        tt.fields.namespaces,
				ctx:               tt.fields.ctx,
				enabledCollectors: tt.fields.enabledCollectors,
				whiteBlackList:    tt.fields.whiteBlackList,
			}
			if got := b.WithEnabledCollectors(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.WithEnabledCollectors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithNamespaces(t *testing.T) {
	type fields struct {
		apiserver         string
		kubeconfig        string
		namespaces        koptions.NamespaceList
		ctx               context.Context
		enabledCollectors []string
		whiteBlackList    whiteBlackLister
	}
	type args struct {
		n koptions.NamespaceList
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "namespace",
			fields: fields{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
			},
			args: args{
				n: []string{"ns1"},
			},
			want: &Builder{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        []string{"ns1"},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				apiserver:         tt.fields.apiserver,
				kubeconfig:        tt.fields.kubeconfig,
				namespaces:        tt.fields.namespaces,
				ctx:               tt.fields.ctx,
				enabledCollectors: tt.fields.enabledCollectors,
				whiteBlackList:    tt.fields.whiteBlackList,
			}
			if got := b.WithNamespaces(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.WithNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithWhiteBlackList(t *testing.T) {
	w, _ := whiteblacklist.New(map[string]struct{}{}, map[string]struct{}{})
	type fields struct {
		apiserver         string
		kubeconfig        string
		namespaces        koptions.NamespaceList
		ctx               context.Context
		enabledCollectors []string
		whiteBlackList    whiteBlackLister
	}
	type args struct {
		l whiteBlackLister
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "whiteBlackList",
			fields: fields{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
				whiteBlackList:    nil,
			},
			args: args{
				l: w,
			},
			want: &Builder{
				apiserver:         "",
				kubeconfig:        "",
				namespaces:        koptions.NamespaceList{},
				ctx:               ctx,
				enabledCollectors: []string{"col1", "col2"},
				whiteBlackList:    w,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				apiserver:         tt.fields.apiserver,
				kubeconfig:        tt.fields.kubeconfig,
				namespaces:        tt.fields.namespaces,
				ctx:               tt.fields.ctx,
				enabledCollectors: tt.fields.enabledCollectors,
				whiteBlackList:    tt.fields.whiteBlackList,
			}
			if got := b.WithWhiteBlackList(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.WithWhiteBlackList() = %v, want %v", got, tt.want)
			}
		})
	}
}
