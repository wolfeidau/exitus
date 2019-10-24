import { resolve } from 'path';
import { Stack, App, StackProps, CfnOutput, Duration } from '@aws-cdk/core'
import { DatabaseInstance, DatabaseInstanceEngine } from '@aws-cdk/aws-rds'
import { ManagedPolicy, ServicePrincipal, Role } from '@aws-cdk/aws-iam'
import { Vpc, InstanceType, InstanceClass, SubnetType, InstanceSize, Port } from '@aws-cdk/aws-ec2';
import { Cluster, ContainerImage, Secret } from '@aws-cdk/aws-ecs';
import { ApplicationLoadBalancedFargateService } from '@aws-cdk/aws-ecs-patterns';
import { HostedZone, ARecord, AddressRecordTarget } from '@aws-cdk/aws-route53';
import { LoadBalancerTarget } from '@aws-cdk/aws-route53-targets';
import { Certificate } from '@aws-cdk/aws-certificatemanager';

const appName = "Exitus"
const stage = process.env.STAGE || 'dev';
const branch = process.env.BRANCH || 'master';
const oauthClientID = process.env.OAUTH_CLIENT_ID || '';
const oauthClientSecret = process.env.OAUTH_CLIENT_SECRET || '';
const openidProviderURL = process.env.OPENID_PROVIDER_URL || '';
const dbPort = 5432;

const domainName = process.env.DOMAIN_NAME || ''
const hostedZoneID = process.env.HOSTED_ZONE_ID || ''
const acmCertificateArn = process.env.ACM_CERTIFICATE_ARN || ''

interface RDSStackProps extends StackProps {
    vpc: Vpc
}
interface ServiceStackProps extends StackProps {
    vpc: Vpc
    db: DatabaseInstance;
}

class VPCStack extends Stack {
    vpc: Vpc;
    constructor(scope: App, id: string, props?: StackProps) {
        super(scope, id, props);
        // Network to run everything in
        this.vpc = new Vpc(this, `${appName}Vpc`, {
            cidr: "10.0.0.0/16",
            maxAzs: 2,
            natGateways: 1,
            subnetConfiguration: [
                {
                    name: `${appName.toLowerCase()}-${stage}-public`,
                    subnetType: SubnetType.PUBLIC,
                },
                {
                    name: `${appName.toLowerCase()}-${stage}-private`,
                    subnetType: SubnetType.PRIVATE,
                },
                {
                    name: `${appName.toLowerCase()}-${stage}-isolated`,
                    subnetType: SubnetType.ISOLATED,
                },
            ]
        });
    }
}

class RDSStack extends Stack {
    db: DatabaseInstance;
    constructor(scope: App, id: string, props?: RDSStackProps) {
        super(scope, id, props);

        this.db = new DatabaseInstance(this, `${appName}DB`, {
            databaseName: appName.toLowerCase(),
            masterUsername: appName.toLowerCase(),
            engine: DatabaseInstanceEngine.POSTGRES,
            instanceClass: InstanceType.of(
                InstanceClass.BURSTABLE3, InstanceSize.SMALL
            ),
            vpc: props!.vpc,
            vpcPlacement: {
                subnetType: SubnetType.PRIVATE
            },
            port: dbPort,
        })

    }
}

class ServiceStack extends Stack {
    constructor(scope: App, id: string, props?: ServiceStackProps) {
        super(scope, id, props);

        const cluster = new Cluster(this, `${appName}Cluster`, { vpc: props!.vpc });

        const executionRole = new Role(this, `${appName}ExecutionRole`, {
            assumedBy: new ServicePrincipal('ecs-tasks.amazonaws.com'),
            managedPolicies: [
                ManagedPolicy.fromAwsManagedPolicyName('service-role/AmazonECSTaskExecutionRolePolicy')
            ]
        })
        const taskRole = new Role(this, `${appName}TaskRole`, {
            assumedBy: new ServicePrincipal('ecs-tasks.amazonaws.com'),
        })

        const certificate = Certificate.fromCertificateArn(this,`${appName}Cert`, acmCertificateArn)

        const containerSvc = new ApplicationLoadBalancedFargateService(this, `${appName}FargateService`, {
            cluster,
            image: ContainerImage.fromAsset(resolve(__dirname, './deploy')),
            environment: {
                OAUTH_CLIENT_ID: oauthClientID,
                OAUTH_CLIENT_SECRET: oauthClientSecret,
                OPENID_PROVIDER_URL: openidProviderURL,
                BRANCH: branch,
                STAGE: stage,
                ADDR: ":80",
            },
            secrets : {
                DB_SECRET: Secret.fromSecretsManager(props!.db.secret!),
            },
            certificate,
            enableLogging: true,
            containerPort: 80,
            executionRole,
            taskRole,
        });

        containerSvc.targetGroup.healthCheck = {
            interval: Duration.seconds(10),
            path: "/healthz",
            timeout: Duration.seconds(5),
        }

        const zone = HostedZone.fromHostedZoneId(this, 'MyZone', hostedZoneID);

        const siteDomain = `${appName.toLowerCase()}-${stage}-${branch}` + '.' + domainName;

        new ARecord(this, `${appName}Record`, {
            zone,
            target: AddressRecordTarget.fromAlias(new LoadBalancerTarget(containerSvc.loadBalancer)),
            recordName: siteDomain,
        })

        containerSvc.service.connections.allowTo(props!.db!, Port.tcp(dbPort))

        new CfnOutput(this, 'LoadBalancerDNS', { value: containerSvc.loadBalancer.loadBalancerDnsName });
        new CfnOutput(this, 'SiteDomain', { value: siteDomain });
    }
}

const app = new App();

const vpcStack = new VPCStack(app, `${appName.toLowerCase()}-vpc-${stage}`, {});

const rdsStack = new RDSStack(app, `${appName.toLowerCase()}-rds-${stage}`, {
    vpc: vpcStack.vpc,
});

new ServiceStack(app, `${appName.toLowerCase()}-service-${stage}`, {
    vpc: vpcStack.vpc,
    db: rdsStack.db,
});

app.synth();
