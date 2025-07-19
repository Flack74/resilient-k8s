-- Insert a test experiment
INSERT INTO experiments (
    id, 
    name, 
    description, 
    type, 
    status, 
    target, 
    parameters, 
    created_at, 
    updated_at, 
    duration
) VALUES (
    'test-experiment-001',
    'Pod Failure Test',
    'Test resilience by killing frontend pods',
    'pod-failure',
    'pending',
    'app=frontend',
    '{"namespace": "default", "percentage": "50"}',
    NOW(),
    NOW(),
    30
);