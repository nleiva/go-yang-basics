go install github.com/openconfig/ygot/generator@latest
generator -path=. \
  -output_file=pkg/base.go \
  -enum_suffix_for_simple_union_enums \
  -package_name=test -generate_fakeroot -fakeroot_name=test \
  -generate_getters \
  -generate_ordered_maps=false \
  -generate_simple_unions \
  base.yang