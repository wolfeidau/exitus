"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const path_1 = require("path");
const core_1 = require("@aws-cdk/core");
const aws_rds_1 = require("@aws-cdk/aws-rds");
const aws_iam_1 = require("@aws-cdk/aws-iam");
const aws_ec2_1 = require("@aws-cdk/aws-ec2");
const aws_ecs_1 = require("@aws-cdk/aws-ecs");
const aws_ecs_patterns_1 = require("@aws-cdk/aws-ecs-patterns");
const appName = "Exitus";
const envName = process.env.ENV_NAME || 'dev';
const dbPort = 5432;
class VPCStack extends core_1.Stack {
    constructor(scope, id, props) {
        super(scope, id, props);
        // Network to run everything in
        this.vpc = new aws_ec2_1.Vpc(this, `${appName}Vpc`, {
            cidr: "10.0.0.0/16",
            maxAzs: 2,
            natGateways: 1,
            subnetConfiguration: [
                {
                    name: `${appName.toLowerCase()}-${envName}-public`,
                    subnetType: aws_ec2_1.SubnetType.PUBLIC,
                },
                {
                    name: `${appName.toLowerCase()}-${envName}-private`,
                    subnetType: aws_ec2_1.SubnetType.PRIVATE,
                },
                {
                    name: `${appName.toLowerCase()}-${envName}-isolated`,
                    subnetType: aws_ec2_1.SubnetType.ISOLATED,
                },
            ]
        });
    }
}
class RDSStack extends core_1.Stack {
    constructor(scope, id, props) {
        super(scope, id, props);
        this.db = new aws_rds_1.DatabaseInstance(this, `${appName}DB`, {
            databaseName: appName.toLowerCase(),
            masterUsername: appName.toLowerCase(),
            engine: aws_rds_1.DatabaseInstanceEngine.POSTGRES,
            instanceClass: aws_ec2_1.InstanceType.of(aws_ec2_1.InstanceClass.BURSTABLE2, aws_ec2_1.InstanceSize.SMALL),
            vpc: props.vpc,
            vpcPlacement: {
                subnetType: aws_ec2_1.SubnetType.PRIVATE
            },
            port: dbPort,
        });
    }
}
class ServiceStack extends core_1.Stack {
    constructor(scope, id, props) {
        super(scope, id, props);
        const cluster = new aws_ecs_1.Cluster(this, `${appName}Cluster`, { vpc: props.vpc });
        const executionRole = new aws_iam_1.Role(this, `${appName}ExecutionRole`, {
            assumedBy: new aws_iam_1.ServicePrincipal('ecs-tasks.amazonaws.com'),
            managedPolicies: [
                aws_iam_1.ManagedPolicy.fromAwsManagedPolicyName('service-role/AmazonECSTaskExecutionRolePolicy')
            ]
        });
        const taskRole = new aws_iam_1.Role(this, `${appName}TaskRole`, {
            assumedBy: new aws_iam_1.ServicePrincipal('ecs-tasks.amazonaws.com'),
        });
        const containerSvc = new aws_ecs_patterns_1.LoadBalancedFargateService(this, `${appName}FargateService`, {
            cluster,
            image: aws_ecs_1.ContainerImage.fromAsset(path_1.resolve(__dirname, './deploy')),
            environment: {
                DB_USERNAME: props.db.secret.secretValueFromJson('username').toString(),
                DB_PASSWORD: props.db.secret.secretValueFromJson('password').toString(),
                DB_HOST: props.db.secret.secretValueFromJson('host').toString(),
                DB_PORT: props.db.secret.secretValueFromJson('port').toString(),
                DB_DATABASE: appName.toLowerCase(),
            },
            executionRole,
            taskRole,
        });
        containerSvc.service.connections.allowTo(props.db, aws_ec2_1.Port.tcp(dbPort));
    }
}
const app = new core_1.App();
const vpcStack = new VPCStack(app, `${appName.toLowerCase()}-vpc-${envName}`, {});
const rdsStack = new RDSStack(app, `${appName.toLowerCase()}-rds-${envName}`, {
    vpc: vpcStack.vpc,
});
new ServiceStack(app, `${appName.toLowerCase()}-service-${envName}`, {
    vpc: vpcStack.vpc,
    db: rdsStack.db,
});
app.synth();
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXhpdHVzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vZXhpdHVzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEsK0JBQStCO0FBQy9CLHdDQUFzRDtBQUN0RCw4Q0FBMkU7QUFDM0UsOENBQXdFO0FBQ3hFLDhDQUFvRztBQUNwRyw4Q0FBMkQ7QUFDM0QsZ0VBQXVFO0FBRXZFLE1BQU0sT0FBTyxHQUFHLFFBQVEsQ0FBQTtBQUN4QixNQUFNLE9BQU8sR0FBRyxPQUFPLENBQUMsR0FBRyxDQUFDLFFBQVEsSUFBSSxLQUFLLENBQUM7QUFDOUMsTUFBTSxNQUFNLEdBQUcsSUFBSSxDQUFDO0FBVXBCLE1BQU0sUUFBUyxTQUFRLFlBQUs7SUFFeEIsWUFBWSxLQUFVLEVBQUUsRUFBVSxFQUFFLEtBQWtCO1FBQ2xELEtBQUssQ0FBQyxLQUFLLEVBQUUsRUFBRSxFQUFFLEtBQUssQ0FBQyxDQUFDO1FBQ3hCLCtCQUErQjtRQUMvQixJQUFJLENBQUMsR0FBRyxHQUFHLElBQUksYUFBRyxDQUFDLElBQUksRUFBRSxHQUFHLE9BQU8sS0FBSyxFQUFFO1lBQ3RDLElBQUksRUFBRSxhQUFhO1lBQ25CLE1BQU0sRUFBRSxDQUFDO1lBQ1QsV0FBVyxFQUFFLENBQUM7WUFDZCxtQkFBbUIsRUFBRTtnQkFDakI7b0JBQ0ksSUFBSSxFQUFFLEdBQUcsT0FBTyxDQUFDLFdBQVcsRUFBRSxJQUFJLE9BQU8sU0FBUztvQkFDbEQsVUFBVSxFQUFFLG9CQUFVLENBQUMsTUFBTTtpQkFDaEM7Z0JBQ0Q7b0JBQ0ksSUFBSSxFQUFFLEdBQUcsT0FBTyxDQUFDLFdBQVcsRUFBRSxJQUFJLE9BQU8sVUFBVTtvQkFDbkQsVUFBVSxFQUFFLG9CQUFVLENBQUMsT0FBTztpQkFDakM7Z0JBQ0Q7b0JBQ0ksSUFBSSxFQUFFLEdBQUcsT0FBTyxDQUFDLFdBQVcsRUFBRSxJQUFJLE9BQU8sV0FBVztvQkFDcEQsVUFBVSxFQUFFLG9CQUFVLENBQUMsUUFBUTtpQkFDbEM7YUFD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QUksVUFBRyxFQUFFLENBQUM7QUFFdEIsTUFBTSxRQUFRLEdBQUcsSUFBSSxRQUFRLENBQUMsR0FBRyxFQUFFLEdBQUcsT0FBTyxDQUFDLFdBQVcsRUFBRSxRQUFRLE9BQU8sRUFBRSxFQUFFLEVBQUUsQ0FBQyxDQUFDO0FBRWxGLE1BQU0sUUFBUSxHQUFHLElBQUksUUFBUSxDQUFDLEdBQUcsRUFBRSxHQUFHLE9BQU8sQ0FBQyxXQUFXLEVBQUUsUUFBUSxPQUFPLEVBQUUsRUFBRTtJQUMxRSxHQUFHLEVBQUUsUUFBUSxDQUFDLEdBQUc7Q0FDcEIsQ0FBQyxDQUFDO0FBRUgsSUFBSSxZQUFZLENBQUMsR0FBRyxFQUFFLEdBQUcsT0FBTyxDQUFDLFdBQVcsRUFBRSxZQUFZLE9BQU8sRUFBRSxFQUFFO0lBQ2pFLEdBQUcsRUFBRSxRQUFRLENBQUMsR0FBRztJQUNqQixFQUFFLEVBQUUsUUFBUSxDQUFDLEVBQUU7Q0FDbEIsQ0FBQyxDQUFDO0FBRUgsR0FBRyxDQUFDLEtBQUssRUFBRSxDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgcmVzb2x2ZSB9IGZyb20gJ3BhdGgnO1xuaW1wb3J0IHsgU3RhY2ssIEFwcCwgU3RhY2tQcm9wcyB9IGZyb20gJ0Bhd3MtY2RrL2NvcmUnXG5pbXBvcnQgeyBEYXRhYmFzZUluc3RhbmNlLCBEYXRhYmFzZUluc3RhbmNlRW5naW5lIH0gZnJvbSAnQGF3cy1jZGsvYXdzLXJkcydcbmltcG9ydCB7IE1hbmFnZWRQb2xpY3ksIFNlcnZpY2VQcmluY2lwYWwsIFJvbGUgfSBmcm9tICdAYXdzLWNkay9hd3MtaWFtJ1xuaW1wb3J0IHsgVnBjLCBJbnN0YW5jZVR5cGUsIEluc3RhbmNlQ2xhc3MsIFN1Ym5ldFR5cGUsIEluc3RhbmNlU2l6ZSwgUG9ydCB9IGZyb20gJ0Bhd3MtY2RrL2F3cy1lYzInO1xuaW1wb3J0IHsgQ2x1c3RlciwgQ29udGFpbmVySW1hZ2UgfSBmcm9tICdAYXdzLWNkay9hd3MtZWNzJztcbmltcG9ydCB7IExvYWRCYWxhbmNlZEZhcmdhdGVTZXJ2aWNlIH0gZnJvbSAnQGF3cy1jZGsvYXdzLWVjcy1wYXR0ZXJucyc7XG5cbmNvbnN0IGFwcE5hbWUgPSBcIkV4aXR1c1wiXG5jb25zdCBlbnZOYW1lID0gcHJvY2Vzcy5lbnYuRU5WX05BTUUgfHwgJ2Rldic7XG5jb25zdCBkYlBvcnQgPSA1NDMyO1xuXG5pbnRlcmZhY2UgUkRTU3RhY2tQcm9wcyBleHRlbmRzIFN0YWNrUHJvcHMge1xuICAgIHZwYzogVnBjXG59XG5pbnRlcmZhY2UgU2VydmljZVN0YWNrUHJvcHMgZXh0ZW5kcyBTdGFja1Byb3BzIHtcbiAgICB2cGM6IFZwY1xuICAgIGRiOiBEYXRhYmFzZUluc3RhbmNlO1xufVxuXG5jbGFzcyBWUENTdGFjayBleHRlbmRzIFN0YWNrIHtcbiAgICB2cGM6IFZwYztcbiAgICBjb25zdHJ1Y3RvcihzY29wZTogQXBwLCBpZDogc3RyaW5nLCBwcm9wcz86IFN0YWNrUHJvcHMpIHtcbiAgICAgICAgc3VwZXIoc2NvcGUsIGlkLCBwcm9wcyk7XG4gICAgICAgIC8vIE5ldHdvcmsgdG8gcnVuIGV2ZXJ5dGhpbmcgaW5cbiAgICAgICAgdGhpcy52cGMgPSBuZXcgVnBjKHRoaXMsIGAke2FwcE5hbWV9VnBjYCwge1xuICAgICAgICAgICAgY2lkcjogXCIxMC4wLjAuMC8xNlwiLFxuICAgICAgICAgICAgbWF4QXpzOiAyLFxuICAgICAgICAgICAgbmF0R2F0ZXdheXM6IDEsXG4gICAgICAgICAgICBzdWJuZXRDb25maWd1cmF0aW9uOiBbXG4gICAgICAgICAgICAgICAge1xuICAgICAgICAgICAgICAgICAgICBuYW1lOiBgJHthcHBOYW1lLnRvTG93ZXJDYXNlKCl9LSR7ZW52TmFtZX0tcHVibGljYCxcbiAgICAgICAgICAgICAgICAgICAgc3VibmV0VHlwZTogU3VibmV0VHlwZS5QVUJMSUMsXG4gICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICB7XG4gICAgICAgICAgICAgICAgICAgIG5hbWU6IGAke2FwcE5hbWUudG9Mb3dlckNhc2UoKX0tJHtlbnZOYW1lfS1wcml2YXRlYCxcbiAgICAgICAgICAgICAgICAgICAgc3VibmV0VHlwZTogU3VibmV0VHlwZS5QUklWQVRFLFxuICAgICAgICAgICAgICAgIH0sXG4gICAgICAgICAgICAgICAge1xuICAgICAgICAgICAgICAgICAgICBuYW1lOiBgJHthcHBOYW1lLnRvTG93ZXJDYXNlKCl9LSR7ZW52TmFtZX0taXNvbGF0ZWRgLFxuICAgICAgICAgICAgICAgICAgICBzdWJuZXRUeXBlOiBTdWJuZXRUeXBlLklTT0xBVEVELFxuICAgICAgICAgICAgICAgIH0sXG4gICAgICAgICAgICBdXG4gICAgICAgIH0pO1xuICAgIH1cbn1cblxuY2xhc3MgUkRTU3RhY2sgZXh0ZW5kcyBTdGFjayB7XG4gICAgZGI6IERhdGFiYXNlSW5zdGFuY2U7XG4gICAgY29uc3RydWN0b3Ioc2NvcGU6IEFwcCwgaWQ6IHN0cmluZywgcHJvcHM/OiBSRFNTdGFja1Byb3BzKSB7XG4gICAgICAgIHN1cGVyKHNjb3BlLCBpZCwgcHJvcHMpO1xuXG4gICAgICAgIHRoaXMuZGIgPSBuZXcgRGF0YWJhc2VJbnN0YW5jZSh0aGlzLCBgJHthcHBOYW1lfURCYCwge1xuICAgICAgICAgICAgZGF0YWJhc2VOYW1lOiBhcHBOYW1lLnRvTG93ZXJDYXNlKCksXG4gICAgICAgICAgICBtYXN0ZXJVc2VybmFtZTogYXBwTmFtZS50b0xvd2VyQ2FzZSgpLFxuICAgICAgICAgICAgZW5naW5lOiBEYXRhYmFzZUluc3RhbmNlRW5naW5lLlBPU1RHUkVTLFxuICAgICAgICAgICAgaW5zdGFuY2VDbGFzczogSW5zdGFuY2VUeXBlLm9mKFxuICAgICAgICAgICAgICAgIEluc3RhbmNlQ2xhc3MuQlVSU1RBQkxFMiwgSW5zdGFuY2VTaXplLlNNQUxMXG4gICAgICAgICAgICApLFxuICAgICAgICAgICAgdnBjOiBwcm9wcyEudnBjLFxuICAgICAgICAgICAgdnBjUGxhY2VtZW50OiB7XG4gICAgICAgICAgICAgICAgc3VibmV0VHlwZTogU3VibmV0VHlwZS5QUklWQVRFXG4gICAgICAgICAgICB9LFxuICAgICAgICAgICAgcG9ydDogZGJQb3J0LFxuICAgICAgICB9KVxuXG4gICAgfVxufVxuXG5jbGFzcyBTZXJ2aWNlU3RhY2sgZXh0ZW5kcyBTdGFjayB7XG4gICAgY29uc3RydWN0b3Ioc2NvcGU6IEFwcCwgaWQ6IHN0cmluZywgcHJvcHM/OiBTZXJ2aWNlU3RhY2tQcm9wcykge1xuICAgICAgICBzdXBlcihzY29wZSwgaWQsIHByb3BzKTtcblxuICAgICAgICBjb25zdCBjbHVzdGVyID0gbmV3IENsdXN0ZXIodGhpcywgYCR7YXBwTmFtZX1DbHVzdGVyYCwgeyB2cGM6IHByb3BzIS52cGMgfSk7XG5cbiAgICAgICAgY29uc3QgZXhlY3V0aW9uUm9sZSA9IG5ldyBSb2xlKHRoaXMsIGAke2FwcE5hbWV9RXhlY3V0aW9uUm9sZWAsIHtcbiAgICAgICAgICAgIGFzc3VtZWRCeTogbmV3IFNlcnZpY2VQcmluY2lwYWwoJ2Vjcy10YXNrcy5hbWF6b25hd3MuY29tJyksXG4gICAgICAgICAgICBtYW5hZ2VkUG9saWNpZXM6IFtcbiAgICAgICAgICAgICAgICBNYW5hZ2VkUG9saWN5LmZyb21Bd3NNYW5hZ2VkUG9saWN5TmFtZSgnc2VydmljZS1yb2xlL0FtYXpvbkVDU1Rhc2tFeGVjdXRpb25Sb2xlUG9saWN5JylcbiAgICAgICAgICAgIF1cbiAgICAgICAgfSlcbiAgICAgICAgY29uc3QgdGFza1JvbGUgPSBuZXcgUm9sZSh0aGlzLCBgJHthcHBOYW1lfVRhc2tSb2xlYCwge1xuICAgICAgICAgICAgYXNzdW1lZEJ5OiBuZXcgU2VydmljZVByaW5jaXBhbCgnZWNzLXRhc2tzLmFtYXpvbmF3cy5jb20nKSxcbiAgICAgICAgfSlcblxuICAgICAgICBjb25zdCBjb250YWluZXJTdmMgPSBuZXcgTG9hZEJhbGFuY2VkRmFyZ2F0ZVNlcnZpY2UodGhpcywgYCR7YXBwTmFtZX1GYXJnYXRlU2VydmljZWAsIHtcbiAgICAgICAgICAgIGNsdXN0ZXIsXG4gICAgICAgICAgICBpbWFnZTogQ29udGFpbmVySW1hZ2UuZnJvbUFzc2V0KHJlc29sdmUoX19kaXJuYW1lLCAnLi9kZXBsb3knKSksXG4gICAgICAgICAgICBlbnZpcm9ubWVudDoge1xuICAgICAgICAgICAgICAgIERCX1VTRVJOQU1FOiBwcm9wcyEuZGIuc2VjcmV0IS5zZWNyZXRWYWx1ZUZyb21Kc29uKCd1c2VybmFtZScpLnRvU3RyaW5nKCksXG4gICAgICAgICAgICAgICAgREJfUEFTU1dPUkQ6IHByb3BzIS5kYi5zZWNyZXQhLnNlY3JldFZhbHVlRnJvbUpzb24oJ3Bhc3N3b3JkJykudG9TdHJpbmcoKSxcbiAgICAgICAgICAgICAgICBEQl9IT1NUOiBwcm9wcyEuZGIuc2VjcmV0IS5zZWNyZXRWYWx1ZUZyb21Kc29uKCdob3N0JykudG9TdHJpbmcoKSxcbiAgICAgICAgICAgICAgICBEQl9QT1JUOiBwcm9wcyEuZGIuc2VjcmV0IS5zZWNyZXRWYWx1ZUZyb21Kc29uKCdwb3J0JykudG9TdHJpbmcoKSxcbiAgICAgICAgICAgICAgICBEQl9EQVRBQkFTRTogYXBwTmFtZS50b0xvd2VyQ2FzZSgpLFxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIGV4ZWN1dGlvblJvbGUsXG4gICAgICAgICAgICB0YXNrUm9sZSxcbiAgICAgICAgfSk7XG5cbiAgICAgICAgY29udGFpbmVyU3ZjLnNlcnZpY2UuY29ubmVjdGlvbnMuYWxsb3dUbyhwcm9wcyEuZGIhLCBQb3J0LnRjcChkYlBvcnQpKVxuICAgIH1cbn1cblxuY29uc3QgYXBwID0gbmV3IEFwcCgpO1xuXG5jb25zdCB2cGNTdGFjayA9IG5ldyBWUENTdGFjayhhcHAsIGAke2FwcE5hbWUudG9Mb3dlckNhc2UoKX0tdnBjLSR7ZW52TmFtZX1gLCB7fSk7XG5cbmNvbnN0IHJkc1N0YWNrID0gbmV3IFJEU1N0YWNrKGFwcCwgYCR7YXBwTmFtZS50b0xvd2VyQ2FzZSgpfS1yZHMtJHtlbnZOYW1lfWAsIHtcbiAgICB2cGM6IHZwY1N0YWNrLnZwYyxcbn0pO1xuXG5uZXcgU2VydmljZVN0YWNrKGFwcCwgYCR7YXBwTmFtZS50b0xvd2VyQ2FzZSgpfS1zZXJ2aWNlLSR7ZW52TmFtZX1gLCB7XG4gICAgdnBjOiB2cGNTdGFjay52cGMsXG4gICAgZGI6IHJkc1N0YWNrLmRiLFxufSk7XG5cbmFwcC5zeW50aCgpO1xuIl19