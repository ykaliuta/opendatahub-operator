IMAGE_OWNER = ykaliuta
IMG_TAG = latest

# just install it in common space, not separately
#ifeq ($(OPERATOR_NAMESPACE),opendatahub-operator-system)
#OPERATOR_NAMESPACE := opendatahub
#endif

# just a short name
ifneq ($(ODH_NS),)
OPERATOR_NAMESPACE := $(ODH_NS)
endif

OPERATOR_NAME = opendatahub-operator

ifeq ($(OPERATOR_NAMESPACE),redhat-ods-operator)
OPERATOR_NAME = rhoai-operator
endif

# if the namespace was not given and it's bundle build for ODH, use openshift-operators namespace
ifeq ($(ODH_NS),)
ifneq ($(OPERATOR_NAME),rhoai-operator)
ifneq ($(findstring bundle,$(MAKECMDGOALS)),)
OPERATOR_NAMESPACE := openshift-operators
endif
endif
endif
#$(info Operator namespace is $(OPERATOR_NAMESPACE))


IMAGE_TAG_BASE := quay.io/$(IMAGE_OWNER)/$(OPERATOR_NAME)
#IMG is constructed lazy with tag base and tag

# allow to override only skip
E2E_SKIP_DELETION = false
E2E_TEST_FLAGS := --skip-deletion=$(E2E_SKIP_DELETION) -timeout 20m

# allow to override only local
USE_LOCAL = false
IMAGE_BUILD_FLAGS := --build-arg USE_LOCAL=$(USE_LOCAL)

bundle-run:
	$(OPERATOR_SDK) run bundle $(BUNDLE_IMG) --namespace $(OPERATOR_NAMESPACE)
