generator -path=. \
  -output_file=pkg/network.go \
  -enum_suffix_for_simple_union_enums \
  -package_name=network -generate_fakeroot -fakeroot_name=device \
  -generate_getters \
  -generate_ordered_maps=false \
  -generate_simple_unions \
  base.yang \
  deviation.yang \
  augment.yang