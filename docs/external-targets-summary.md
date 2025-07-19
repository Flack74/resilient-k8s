# External Targets Support

## Overview

We've added support for external targets to the Chaos Engineering Platform. This allows users to run chaos experiments against systems outside of their Kubernetes cluster, such as:

- External APIs or web services
- On-premises infrastructure
- Cloud services in different environments
- Third-party services that support chaos testing

## Implementation Details

1. **ExperimentController Interface**: Created a common interface for all experiment controllers
2. **K8sExperimentController**: Renamed the original controller to handle Kubernetes-specific experiments
3. **ExternalExperimentController**: Added a new controller for external targets
4. **Target Type Detection**: Updated the operator to detect external targets and use the appropriate controller
5. **Documentation**: Added comprehensive documentation on using external targets

## How to Use

See the [External Targets Guide](external-targets.md) for detailed instructions on using external targets.

## Example

We've provided a complete example in the `examples/external-target` directory:

- `external-service.go`: A simple external service that can receive chaos experiments
- `test-external-target.sh`: A script to test the external target functionality
- `README.md`: Documentation for the example

## API Changes

When creating an experiment for an external target, include these parameters:

```json
{
  "parameters": {
    "target_type": "external",
    "endpoint": "/inject-failure",
    "auth_token": "your-auth-token",
    "type": "service-failure",
    "cleanup_endpoint": "/reset"
  }
}
```

## Security Considerations

- Use HTTPS for all external target communications
- Implement proper authentication for external targets
- Consider network security implications when targeting external systems
- Ensure you have permission to run chaos experiments against the target system