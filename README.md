# Simple Postgres Database Operator

The Postgres Database Operator is a Kubernetes Operator designed to simplify the deploying and monitoring of PostgreSQL database deployments within a Kubernetes environment. By leveraging the power of Kubernetes Operators, this project automates tasks such as provisioning, scaling, and maintaining the desired state of PostgreSQL instances.

## Installation

### Prerequisites
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [operator-sdk](https://sdk.operatorframework.io/docs/installation/)
- Minikube for the K8s Cluster
- Go 1.21
- Docker

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/devansh/database-op.git
   cd database-op
   kubectl apply -f config/crd/bases/database.devansh.com_postgres.yaml
   make run  (if you are running locally)
   kubectl apply -f config/samples/database_v1alpha1_postgres.yaml
   ```

### Usage
- Ensure that PostgreSQL instances are running and maintain their state within the specified namespace.
- Scale PostgreSQL deployments based on predefined hourly time intervals.



### Extra:
```
operator-sdk init --domain devansh.com --repo github.com/devansh/database-op
operator-sdk create api --group database --version v1alpha1 --kind Postgres --resource --controller
```