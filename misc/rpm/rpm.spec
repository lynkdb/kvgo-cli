%define app_home /usr/local/bin
%define app_user kvgo

Name: kvgo-cli
Version: __version__
Release: __release__%{?dist}
Vendor:  lynkdb.com
Summary: lynkdb key-value database server
License: Apache 2

Source0: %{name}-__version__.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}

%description

%prep
%setup -q -n %{name}-%{version}
%build

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{app_home}

install -m 755 bin/kvgo-cli %{buildroot}%{app_home}/kvgo-cli

%clean
rm -rf %{buildroot}

%pre

%post

%preun

%postun

%files
%{app_home}/kvgo-cli

