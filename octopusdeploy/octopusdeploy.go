package octopusdeploy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dghubble/sling"
)

// Client is an OctopusDeploy for making Octopus API requests.
type Client struct {
	sling                          *sling.Sling
	Accounts                       *accountService
	ActionTemplates                *actionTemplateService
	APIKeys                        *apiKeyService
	Artifacts                      *artifactService
	Authentication                 *authenticationService
	AzureDevOpsConnectivityCheck   *azureDevOpsConnectivityCheckService
	AzureEnvironments              *azureEnvironmentService
	BuildInformation               *buildInformationService
	CertificateConfiguration       *certificateConfigurationService
	Certificates                   *certificateService
	Channels                       *channelService
	CloudTemplate                  *cloudTemplateService
	CommunityActionTemplates       *communityActionTemplateService
	Configuration                  *configurationService
	DashboardConfigurations        *dashboardConfigurationService
	Dashboards                     *dashboardService
	DeploymentProcesses            *deploymentProcessService
	Deployments                    *deploymentService
	DynamicExtensions              *dynamicExtensionService
	Environments                   *environmentService
	Events                         *eventService
	ExternalSecurityGroupProviders *externalSecurityGroupProviderService
	FeaturesConfiguration          *featuresConfigurationService
	Feeds                          *feedService
	Interruptions                  *interruptionService
	Invitations                    *invitationService
	IssueTrackers                  *issueTrackerService
	JiraIntegration                *jiraIntegrationService
	LetsEncryptConfiguration       *letsEncryptConfigurationService
	LibraryVariableSets            *libraryVariableSetService
	Licenses                       *licenseService
	Lifecycles                     *lifecycleService
	MachinePolicies                *machinePolicyService
	MachineRoles                   *machineRoleService
	Machines                       *machineService
	MaintenanceConfiguration       *maintenanceConfigurationService
	Migrations                     *migrationService
	OctopusPackageMetadata         *octopusPackageMetadataService
	OctopusServerNodes             *octopusServerNodeService
	Packages                       *packageService
	PackageMetadata                *packageMetadataService
	PerformanceConfiguration       *performanceConfigurationService
	Permissions                    *permissionService
	ProjectGroups                  *projectGroupService
	Projects                       *projectService
	ProjectTriggers                *projectTriggerService
	Proxies                        *proxyService
	Releases                       *releaseService
	Reporting                      *reportingService
	RunbookProcesses               *runbookProcessService
	RunbookRuns                    *runbookRunService
	Runbooks                       *runbookService
	RunbookSnapshots               *runbookSnapshotService
	Root                           *rootService
	ScheduledProjectTriggers       *scheduledProjectTriggerService
	Scheduler                      *schedulerService
	ScopedUserRoles                *scopedUserRoleService
	ServerConfiguration            *serverConfigurationService
	ServerStatus                   *serverStatuService
	SMTPConfiguration              *smtpConfigurationService
	Spaces                         *spaceService
	Subscriptions                  *subscriptionService
	TagSets                        *tagSetService
	Tasks                          *taskService
	TeamMembership                 *teamMembershipService
	Teams                          *teamService
	Tenants                        *tenantService
	TenantVariables                *tenantVariableService
	UpgradeConfiguration           *upgradeConfigurationService
	UserOnboarding                 *userOnboardingService
	UserRoles                      *userRoleService
	Users                          *userService
	Variables                      *variableService
	WorkerPools                    *workerPoolService
	Workers                        *workerService
	WorkerToolsLatestImages        *workerToolsLatestImageService
}

// NewClient returns a new Octopus API client. If a nil client is provided, a
// new http.Client will be used.
func NewClient(httpClient *http.Client, apiURL *url.URL, apiKey string, spaceID string) (*Client, error) {
	if apiURL == nil {
		return nil, createInvalidParameterError(clientNewClient, ParameterOctopusURL)
	}

	if isEmpty(apiKey) {
		return nil, createInvalidParameterError(clientNewClient, ParameterAPIKey)
	}

	if !isAPIKey(apiKey) {
		return nil, createInvalidParameterError(clientNewClient, ParameterAPIKey)
	}

	baseURLWithAPI := strings.TrimRight(apiURL.String(), "/")
	baseURLWithAPI = fmt.Sprintf("%s/api", baseURLWithAPI)

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// fetch root resource and process paths
	base := sling.New().Client(httpClient).Base(baseURLWithAPI).Set(clientAPIKeyHTTPHeader, apiKey)
	base.Set("User-Agent", "go-octopusdeploy")
	rootService := newRootService(base, baseURLWithAPI)

	root, err := rootService.Get()
	if err != nil {
		return nil, err
	}

	rootPath := root.Links[linkSelf]
	apiKeysPath := "/api/users"
	dynamicExtensionsPath := "/api/dynamic-extensions"
	jiraIntegrationPath := "/api/jiraintegration"
	licensesPath := "/api/licenses"
	migrationsPath := "/api/migrations"
	reportingPath := "/api/reporting"

	accountsPath := root.Links[linkAccounts]
	actionTemplatesLogo := root.Links[linkActionTemplateLogo]
	actionTemplatesPath := root.Links[linkActionTemplates]
	actionTemplatesCategories := root.Links[linkActionTemplatesCategories]
	actionTemplatesSearch := root.Links[linkActionTemplatesSearch]
	actionTemplateVersionedLogo := root.Links[linkActionTemplateVersionedLogo]
	artifactsPath := root.Links[linkArtifacts]
	authenticateOctopusIDPath := root.Links[linkAuthenticateOctopusID]
	authenticationPath := root.Links[linkAuthentication]
	azureDevOpsConnectivityCheckPath := root.Links[linkAzureDevOpsConnectivityCheck]
	azureEnvironmentsPath := root.Links[linkAzureEnvironments]
	buildInformationPath := root.Links[linkBuildInformation]
	buildInformationBulkPath := root.Links[linkBuildInformationBulk]
	builtInFeedStatsPath := root.Links[linkBuiltInFeedStats]
	certificateConfigurationPath := root.Links[linkCertificateConfiguration]
	certificatesPath := root.Links[linkCertificates]
	channelsPath := root.Links[linkChannels]
	cloudTemplatePath := root.Links[linkCloudTemplate]
	communityActionTemplatesPath := root.Links[linkCommunityActionTemplates]
	configurationPath := root.Links[linkConfiguration]
	currentLicensePath := root.Links[linkCurrentLicense]
	currentLicenseStatusPath := root.Links[linkCurrentLicenseStatus]
	currentUserPath := root.Links[linkCurrentUser]
	dashboardPath := root.Links[linkDashboard]
	dashboardConfigurationPath := root.Links[linkDashboardConfiguration]
	dashboardDynamicPath := root.Links[linkDashboardDynamic]
	deploymentProcessesPath := root.Links[linkDeploymentProcesses]
	deploymentsPath := root.Links[linkDeployments]
	discoverMachinePath := root.Links[linkDiscoverMachine]
	discoverWorkerPath := root.Links[linkDiscoverWorker]
	dynamicExtensionsFeaturesMetadataPath := root.Links[linkDynamicExtensionsFeaturesMetadata]
	dynamicExtensionsFeaturesValuesPath := root.Links[linkDynamicExtensionsFeaturesValues]
	dynamicExtensionsScriptsPath := root.Links[linkDynamicExtensionsScripts]
	environmentsPath := root.Links[linkEnvironments]
	environmentSortOrderPath := root.Links[linkEnvironmentSortOrder]
	environmentsSummaryPath := root.Links[linkEnvironmentsSummary]
	eventAgentsPath := root.Links[linkEventAgents]
	eventCategoriesPath := root.Links[linkEventCategories]
	eventDocumentTypesPath := root.Links[linkEventDocumentTypes]
	eventGroupsPath := root.Links[linkEventGroups]
	eventsPath := root.Links[linkEvents]
	extensionStatsPath := root.Links[linkExtensionStats]
	externalSecurityGroupProvidersPath := root.Links[linkExternalSecurityGroupProviders]
	externalUserSearchPath := root.Links[linkExternalUserSearch]
	featuresConfigurationPath := root.Links[linkFeaturesConfiguration]
	feedsPath := root.Links[linkFeeds]
	interruptionsPath := root.Links[linkInterruptions]
	invitationsPath := root.Links[linkInvitations]
	issueTrackersPath := root.Links[linkIssueTrackers]
	jiraConnectAppCredentialsTestPath := root.Links[linkJiraConnectAppCredentialsTest]
	jiraCredentialsTestPath := root.Links[linkJiraCredentialsTest]
	letsEncryptConfigurationPath := root.Links[linkLetsEncryptConfiguration]
	libraryVariablesPath := root.Links[linkLibraryVariables]
	lifecyclesPath := root.Links[linkLifecycles]
	loginInitiatedPath := root.Links[linkLoginInitiated]
	machineOperatingSystemsPath := root.Links[linkMachineOperatingSystems]
	machinePoliciesPath := root.Links[linkMachinePolicies]
	machinePolicyTemplatePath := root.Links[linkMachinePolicyTemplate]
	machineRolesPath := root.Links[linkMachineRoles]
	machinesPath := root.Links[linkMachines]
	machineShellsPath := root.Links[linkMachineShells]
	maintenanceConfigurationPath := root.Links[linkMaintenanceConfiguration]
	migrationsImportPath := root.Links[linkMigrationsImport]
	migrationsPartialExportPath := root.Links[linkMigrationsPartialExport]
	octopusServerClusterSummaryPath := root.Links[linkOctopusServerClusterSummary]
	octopusServerNodesPath := root.Links[linkOctopusServerNodes]
	packageDeltaSignaturePath := root.Links[linkPackageDeltaSignature]
	packageDeltaUploadPath := root.Links[linkPackageDeltaUpload]
	packageMetadataPath := root.Links[linkPackageMetadata]
	packageNotesListPath := root.Links[linkPackageNotesList]
	packagesPath := root.Links[linkPackages]
	packagesBulkPath := root.Links[linkPackagesBulk]
	packageUploadPath := root.Links[linkPackageUpload]
	performanceConfigurationPath := root.Links[linkPerformanceConfiguration]
	permissionsPath := root.Links[linkPermissions]
	projectGroupsPath := root.Links[linkProjectGroups]
	projectPulsePath := root.Links[linkProjectPulse]
	projectsPath := root.Links[linkProjects]
	projectsExperimentalSummariesPath := root.Links[linkProjectsExperimentalSummaries]
	projectTriggersPath := root.Links[linkProjectTriggers]
	proxiesPath := root.Links[linkProxies]
	registerPath := root.Links[linkRegister]
	releasesPath := root.Links[linkReleases]
	reportingDeploymentsCountedByWeekPath := root.Links[linkReportingDeploymentsCountedByWeek]
	runbookProcessesPath := root.Links[linkRunbookProcesses]
	runbookRunsPath := root.Links[linkRunbookRuns]
	runbooksPath := root.Links[linkRunbooks]
	runbookSnapshotsPath := root.Links[linkRunbookSnapshots]
	scheduledProjectTriggersPath := root.Links[linkScheduledProjectTriggers]
	schedulerPath := root.Links[linkScheduler]
	scopedUserRolesPath := root.Links[linkScopedUserRoles]
	serverConfigurationPath := root.Links[linkServerConfiguration]
	serverConfigurationSettingsPath := root.Links[linkServerConfigurationSettings]
	serverHealthStatusPath := root.Links[linkServerHealthStatus]
	serverStatusPath := root.Links[linkServerStatus]
	signInPath := root.Links[linkSignIn]
	signOutPath := root.Links[linkSignOut]
	smtpConfigurationPath := root.Links[linkSMTPConfiguration]
	smtpIsConfiguredPath := root.Links[linkSMTPIsConfigured]
	spaceHomePath := root.Links[linkSpaceHome]
	spacesPath := root.Links[linkSpaces]
	subscriptionsPath := root.Links[linkSubscriptions]
	tagSetsPath := root.Links[linkTagSets]
	tagSetSortOrderPath := root.Links[linkTagSetSortOrder]
	tasksPath := root.Links[linkTasks]
	taskTypesPath := root.Links[linkTaskTypes]
	teamMembershipPath := root.Links[linkTeamMembership]
	teamMembershipPreviewTeamPath := root.Links[linkTeamMembershipPreviewTeam]
	teamsPath := root.Links[linkTeams]
	tenantsPath := root.Links[linkTenants]
	tenantsMissingVariablesPath := root.Links[linkTenantsMissingVariables]
	tenantsStatusPath := root.Links[linkTenantsStatus]
	tenantTagTestPath := root.Links[linkTenantTagTest]
	tenantVariablesPath := root.Links[linkTenantVariables]
	timezonesPath := root.Links[linkTimezones]
	upgradeConfigurationPath := root.Links[linkUpgradeConfiguration]
	userAuthenticationPath := root.Links[linkUserAuthentication]
	userIdentityMetadataPath := root.Links[linkUserIdentityMetadata]
	userOnboardingPath := root.Links[linkUserOnboarding]
	userRolesPath := root.Links[linkUserRoles]
	usersPath := root.Links[linkUsers]
	variableNamesPath := root.Links[linkVariableNames]
	variablePreviewPath := root.Links[linkVariablePreview]
	variablesPath := root.Links[linkVariables]
	versionControlClearCachePath := root.Links[linkVersionControlClearCache]
	versionRuleTestPath := root.Links[linkVersionRuleTest]
	workerOperatingSystemsPath := root.Links[linkWorkerOperatingSystems]
	workerPoolsPath := root.Links[linkWorkerPools]
	workerPoolsDynamicWorkerTypesPath := root.Links[linkWorkerPoolsDynamicWorkerTypes]
	workerPoolsSortOrderPath := root.Links[linkWorkerPoolsSortOrder]
	workerPoolsSummaryPath := root.Links[linkWorkerPoolsSummary]
	workerPoolsSupportedTypesPath := root.Links[linkWorkerPoolsSupportedTypes]
	workersPath := root.Links[linkWorkers]
	workerShellsPath := root.Links[linkWorkerShells]
	workerToolsLatestImagesPath := root.Links[linkWorkerToolsLatestImages]

	if !isEmpty(spaceID) {
		baseURLWithAPI = fmt.Sprintf("%s/%s", baseURLWithAPI, spaceID)
		base = sling.New().Client(httpClient).Base(baseURLWithAPI).Set(clientAPIKeyHTTPHeader, apiKey)
		base.Set("User-Agent", "go-octopusdeploy")
		rootService = newRootService(base, baseURLWithAPI)
		root, err = rootService.Get()

		if err != nil {
			if err == ErrItemNotFound {
				return nil, fmt.Errorf("the space ID (%s) cannot be found", spaceID)
			}
			return nil, err
		}

		if !isEmpty(root.Links[linkAccounts]) {
			accountsPath = root.Links[linkAccounts]
		}

		if !isEmpty(root.Links[linkActionTemplateLogo]) {
			actionTemplatesLogo = root.Links[actionTemplatesLogo]
		}

		if !isEmpty(root.Links[linkActionTemplates]) {
			actionTemplatesPath = root.Links[linkActionTemplates]
		}

		if !isEmpty(root.Links[linkActionTemplatesCategories]) {
			actionTemplatesCategories = root.Links[linkActionTemplatesCategories]
		}

		if !isEmpty(root.Links[linkActionTemplatesSearch]) {
			actionTemplatesSearch = root.Links[linkActionTemplatesSearch]
		}

		if !isEmpty(root.Links[linkActionTemplateVersionedLogo]) {
			actionTemplateVersionedLogo = root.Links[linkActionTemplateVersionedLogo]
		}

		if !isEmpty(root.Links[linkArtifacts]) {
			artifactsPath = root.Links[linkArtifacts]
		}

		if !isEmpty(root.Links[linkAuthenticateOctopusID]) {
			authenticateOctopusIDPath = root.Links[linkAuthenticateOctopusID]
		}

		if !isEmpty(root.Links[linkAuthentication]) {
			authenticationPath = root.Links[linkAuthentication]
		}

		if !isEmpty(root.Links[linkAzureDevOpsConnectivityCheck]) {
			azureDevOpsConnectivityCheckPath = root.Links[linkAzureDevOpsConnectivityCheck]
		}

		if !isEmpty(root.Links[linkAzureEnvironments]) {
			azureEnvironmentsPath = root.Links[linkAzureEnvironments]
		}

		if !isEmpty(root.Links[linkBuildInformation]) {
			buildInformationPath = root.Links[linkBuildInformation]
		}

		if !isEmpty(root.Links[linkBuildInformationBulk]) {
			buildInformationBulkPath = root.Links[linkBuildInformationBulk]
		}

		if !isEmpty(root.Links[linkBuiltInFeedStats]) {
			builtInFeedStatsPath = root.Links[linkBuiltInFeedStats]
		}

		if !isEmpty(root.Links[linkCertificateConfiguration]) {
			certificateConfigurationPath = root.Links[linkCertificateConfiguration]
		}

		if !isEmpty(root.Links[linkCertificates]) {
			certificatesPath = root.Links[linkCertificates]
		}

		if !isEmpty(root.Links[linkChannels]) {
			channelsPath = root.Links[linkChannels]
		}

		if !isEmpty(root.Links[linkCloudTemplate]) {
			cloudTemplatePath = root.Links[linkCloudTemplate]
		}

		if !isEmpty(root.Links[linkCommunityActionTemplates]) {
			communityActionTemplatesPath = root.Links[linkCommunityActionTemplates]
		}

		if !isEmpty(root.Links[linkConfiguration]) {
			configurationPath = root.Links[linkConfiguration]
		}

		if !isEmpty(root.Links[linkCurrentLicense]) {
			currentLicensePath = root.Links[linkCurrentLicense]
		}

		if !isEmpty(root.Links[linkCurrentLicenseStatus]) {
			currentLicenseStatusPath = root.Links[linkCurrentLicenseStatus]
		}

		if !isEmpty(root.Links[linkCurrentUser]) {
			currentUserPath = root.Links[linkCurrentUser]
		}

		if !isEmpty(root.Links[linkDashboard]) {
			dashboardPath = root.Links[linkDashboard]
		}

		if !isEmpty(root.Links[linkDashboardConfiguration]) {
			dashboardConfigurationPath = root.Links[linkDashboardConfiguration]
		}

		if !isEmpty(root.Links[linkDashboardDynamic]) {
			dashboardDynamicPath = root.Links[linkDashboardDynamic]
		}

		if !isEmpty(root.Links[linkDeploymentProcesses]) {
			deploymentProcessesPath = root.Links[linkDeploymentProcesses]
		}

		if !isEmpty(root.Links[linkDeployments]) {
			deploymentsPath = root.Links[linkDeployments]
		}

		if !isEmpty(root.Links[linkDiscoverMachine]) {
			discoverMachinePath = root.Links[linkDiscoverMachine]
		}

		if !isEmpty(root.Links[linkDiscoverWorker]) {
			discoverWorkerPath = root.Links[linkDiscoverWorker]
		}

		if !isEmpty(root.Links[linkDynamicExtensionsFeaturesMetadata]) {
			dynamicExtensionsFeaturesMetadataPath = root.Links[linkDynamicExtensionsFeaturesMetadata]
		}

		if !isEmpty(root.Links[linkDynamicExtensionsFeaturesValues]) {
			dynamicExtensionsFeaturesValuesPath = root.Links[linkDynamicExtensionsFeaturesValues]
		}

		if !isEmpty(root.Links[linkDynamicExtensionsScripts]) {
			dynamicExtensionsScriptsPath = root.Links[linkDynamicExtensionsScripts]
		}

		if !isEmpty(root.Links[linkEnvironments]) {
			environmentsPath = root.Links[linkEnvironments]
		}

		if !isEmpty(root.Links[linkEnvironmentSortOrder]) {
			environmentSortOrderPath = root.Links[linkEnvironmentSortOrder]
		}

		if !isEmpty(root.Links[linkEnvironmentsSummary]) {
			environmentsSummaryPath = root.Links[linkEnvironmentsSummary]
		}

		if !isEmpty(root.Links[linkEventAgents]) {
			eventAgentsPath = root.Links[linkEventAgents]
		}

		if !isEmpty(root.Links[linkEventCategories]) {
			eventCategoriesPath = root.Links[linkEventCategories]
		}

		if !isEmpty(root.Links[linkEventDocumentTypes]) {
			eventDocumentTypesPath = root.Links[linkEventDocumentTypes]
		}

		if !isEmpty(root.Links[linkEventGroups]) {
			eventGroupsPath = root.Links[linkEventGroups]
		}

		if !isEmpty(root.Links[linkEvents]) {
			eventsPath = root.Links[linkEvents]
		}

		if !isEmpty(root.Links[linkExtensionStats]) {
			extensionStatsPath = root.Links[linkExtensionStats]
		}

		if !isEmpty(root.Links[linkExternalSecurityGroupProviders]) {
			externalSecurityGroupProvidersPath = root.Links[linkExternalSecurityGroupProviders]
		}

		if !isEmpty(root.Links[linkExternalUserSearch]) {
			externalUserSearchPath = root.Links[linkExternalUserSearch]
		}

		if !isEmpty(root.Links[linkFeaturesConfiguration]) {
			featuresConfigurationPath = root.Links[linkFeaturesConfiguration]
		}

		if !isEmpty(root.Links[linkFeeds]) {
			feedsPath = root.Links[linkFeeds]
		}

		if !isEmpty(root.Links[linkInterruptions]) {
			interruptionsPath = root.Links[linkInterruptions]
		}

		if !isEmpty(root.Links[linkInvitations]) {
			invitationsPath = root.Links[linkInvitations]
		}

		if !isEmpty(root.Links[linkIssueTrackers]) {
			issueTrackersPath = root.Links[linkIssueTrackers]
		}

		if !isEmpty(root.Links[linkJiraConnectAppCredentialsTest]) {
			jiraConnectAppCredentialsTestPath = root.Links[linkJiraConnectAppCredentialsTest]
		}

		if !isEmpty(root.Links[linkJiraCredentialsTest]) {
			jiraCredentialsTestPath = root.Links[linkJiraCredentialsTest]
		}

		if !isEmpty(root.Links[linkLetsEncryptConfiguration]) {
			letsEncryptConfigurationPath = root.Links[linkLetsEncryptConfiguration]
		}

		if !isEmpty(root.Links[linkLibraryVariables]) {
			libraryVariablesPath = root.Links[linkLibraryVariables]
		}

		if !isEmpty(root.Links[linkLifecycles]) {
			lifecyclesPath = root.Links[linkLifecycles]
		}

		if !isEmpty(root.Links[linkLoginInitiated]) {
			loginInitiatedPath = root.Links[linkLoginInitiated]
		}

		if !isEmpty(root.Links[linkMachineOperatingSystems]) {
			machineOperatingSystemsPath = root.Links[linkMachineOperatingSystems]
		}

		if !isEmpty(root.Links[linkMachinePolicies]) {
			machinePoliciesPath = root.Links[linkMachinePolicies]
		}

		if !isEmpty(root.Links[linkMachinePolicyTemplate]) {
			machinePolicyTemplatePath = root.Links[linkMachinePolicyTemplate]
		}

		if !isEmpty(root.Links[linkMachineRoles]) {
			machineRolesPath = root.Links[linkMachineRoles]
		}

		if !isEmpty(root.Links[linkMachines]) {
			machinesPath = root.Links[linkMachines]
		}

		if !isEmpty(root.Links[linkMachineShells]) {
			machineShellsPath = root.Links[linkMachineShells]
		}

		if !isEmpty(root.Links[linkMaintenanceConfiguration]) {
			maintenanceConfigurationPath = root.Links[linkMaintenanceConfiguration]
		}

		if !isEmpty(root.Links[linkMigrationsImport]) {
			migrationsImportPath = root.Links[linkMigrationsImport]
		}

		if !isEmpty(root.Links[linkMigrationsPartialExport]) {
			migrationsPartialExportPath = root.Links[linkMigrationsPartialExport]
		}

		if !isEmpty(root.Links[linkOctopusServerClusterSummary]) {
			octopusServerClusterSummaryPath = root.Links[linkOctopusServerClusterSummary]
		}

		if !isEmpty(root.Links[linkOctopusServerNodes]) {
			octopusServerNodesPath = root.Links[linkOctopusServerNodes]
		}

		if !isEmpty(root.Links[linkPackageDeltaSignature]) {
			packageDeltaSignaturePath = root.Links[linkPackageDeltaSignature]
		}

		if !isEmpty(root.Links[linkProjects]) {
			projectsPath = root.Links[linkProjects]
		}

		if !isEmpty(root.Links[linkProjectGroups]) {
			projectGroupsPath = root.Links[linkProjectGroups]
		}

		if !isEmpty(root.Links[linkProjectTriggers]) {
			projectTriggersPath = root.Links[linkProjectTriggers]
		}

		if !isEmpty(root.Links[linkSelf]) {
			rootPath = root.Links[linkSelf]
		}

		if !isEmpty(root.Links[linkSpaces]) {
			spacesPath = root.Links[linkSpaces]
		}

		if !isEmpty(root.Links[linkTagSets]) {
			tagSetsPath = root.Links[linkTagSets]
		}

		if !isEmpty(root.Links[linkTenants]) {
			tenantsPath = root.Links[linkTenants]
		}

		if !isEmpty(root.Links[linkUsers]) {
			usersPath = root.Links[linkUsers]
		}

		if !isEmpty(root.Links[linkVariables]) {
			variablesPath = root.Links[linkVariables]
		}

		if !isEmpty(root.Links[linkWorkerPools]) {
			workerPoolsPath = root.Links[linkWorkerPools]
		}

		if !isEmpty(root.Links[linkWorkers]) {
			workersPath = root.Links[linkWorkers]
		}
	}

	return &Client{
		sling:                          base,
		Accounts:                       newAccountService(base, accountsPath),
		ActionTemplates:                newActionTemplateService(base, actionTemplatesPath, actionTemplatesCategories, actionTemplatesLogo, actionTemplatesSearch, actionTemplateVersionedLogo),
		APIKeys:                        newAPIKeyService(base, apiKeysPath),
		Artifacts:                      newArtifactService(base, artifactsPath),
		Authentication:                 newAuthenticationService(base, authenticationPath, loginInitiatedPath),
		AzureDevOpsConnectivityCheck:   newAzureDevOpsConnectivityCheckService(base, azureDevOpsConnectivityCheckPath),
		AzureEnvironments:              newAzureEnvironmentService(base, azureEnvironmentsPath),
		BuildInformation:               newBuildInformationService(base, buildInformationPath, buildInformationBulkPath),
		CertificateConfiguration:       newCertificateConfigurationService(base, certificateConfigurationPath),
		Certificates:                   newCertificateService(base, certificatesPath),
		Channels:                       newChannelService(base, channelsPath, versionRuleTestPath),
		CloudTemplate:                  newCloudTemplateService(base, cloudTemplatePath),
		CommunityActionTemplates:       newCommunityActionTemplateService(base, communityActionTemplatesPath),
		Configuration:                  newConfigurationService(base, configurationPath, versionControlClearCachePath),
		DashboardConfigurations:        newDashboardConfigurationService(base, dashboardConfigurationPath),
		Dashboards:                     newDashboardService(base, dashboardPath, dashboardDynamicPath),
		DeploymentProcesses:            newDeploymentProcessService(base, deploymentProcessesPath),
		Deployments:                    newDeploymentService(base, deploymentsPath),
		DynamicExtensions:              newDynamicExtensionService(base, dynamicExtensionsPath, dynamicExtensionsFeaturesMetadataPath, dynamicExtensionsFeaturesValuesPath, dynamicExtensionsScriptsPath),
		Environments:                   newEnvironmentService(base, environmentsPath, environmentSortOrderPath, environmentsSummaryPath),
		Events:                         newEventService(base, eventsPath, eventAgentsPath, eventCategoriesPath, eventDocumentTypesPath, eventGroupsPath),
		ExternalSecurityGroupProviders: newExternalSecurityGroupProviderService(base, externalSecurityGroupProvidersPath),
		FeaturesConfiguration:          newFeaturesConfigurationService(base, featuresConfigurationPath),
		Feeds:                          newFeedService(base, feedsPath, builtInFeedStatsPath),
		Interruptions:                  newInterruptionService(base, interruptionsPath),
		Invitations:                    newInvitationService(base, invitationsPath),
		IssueTrackers:                  newIssueTrackerService(base, issueTrackersPath),
		JiraIntegration:                newJiraIntegrationService(base, jiraIntegrationPath, jiraConnectAppCredentialsTestPath, jiraCredentialsTestPath),
		LetsEncryptConfiguration:       newLetsEncryptConfigurationService(base, letsEncryptConfigurationPath),
		LibraryVariableSets:            newLibraryVariableSetService(base, libraryVariablesPath),
		Licenses:                       newLicenseService(base, licensesPath, currentLicensePath, currentLicenseStatusPath),
		Lifecycles:                     newLifecycleService(base, lifecyclesPath),
		MachinePolicies:                newMachinePolicyService(base, machinePoliciesPath, machinePolicyTemplatePath),
		MachineRoles:                   newMachineRoleService(base, machineRolesPath),
		Machines:                       newMachineService(base, machinesPath, discoverMachinePath, machineOperatingSystemsPath, machineShellsPath),
		MaintenanceConfiguration:       newMaintenanceConfigurationService(base, maintenanceConfigurationPath),
		Migrations:                     newMigrationService(base, migrationsPath, migrationsImportPath, migrationsPartialExportPath),
		OctopusServerNodes:             newOctopusServerNodeService(base, octopusServerNodesPath, octopusServerClusterSummaryPath),
		Packages:                       newPackageService(base, packagesPath, packageDeltaSignaturePath, packageDeltaUploadPath, packageNotesListPath, packagesBulkPath, packageUploadPath),
		PackageMetadata:                newPackageMetadataService(base, packageMetadataPath),
		PerformanceConfiguration:       newPerformanceConfigurationService(base, performanceConfigurationPath),
		Permissions:                    newPermissionService(base, permissionsPath),
		ProjectGroups:                  newProjectGroupService(base, projectGroupsPath),
		Projects:                       newProjectService(base, projectsPath, projectPulsePath, projectsExperimentalSummariesPath),
		ProjectTriggers:                newProjectTriggerService(base, projectTriggersPath),
		Proxies:                        newProxyService(base, proxiesPath),
		Releases:                       newReleaseService(base, releasesPath),
		Reporting:                      newReportingService(base, reportingPath, reportingDeploymentsCountedByWeekPath),
		RunbookProcesses:               newRunbookProcessService(base, runbookProcessesPath),
		RunbookRuns:                    newRunbookRunService(base, runbookRunsPath),
		Runbooks:                       newRunbookService(base, runbooksPath),
		RunbookSnapshots:               newRunbookSnapshotService(base, runbookSnapshotsPath),
		Root:                           newRootService(base, rootPath),
		Scheduler:                      newSchedulerService(base, schedulerPath),
		ScheduledProjectTriggers:       newScheduledProjectTriggerService(base, scheduledProjectTriggersPath),
		ScopedUserRoles:                newScopedUserRoleService(base, scopedUserRolesPath),
		ServerConfiguration:            newServerConfigurationService(base, serverConfigurationPath, serverConfigurationSettingsPath),
		ServerStatus:                   newServerStatuService(base, serverStatusPath, extensionStatsPath, serverHealthStatusPath, timezonesPath),
		SMTPConfiguration:              newSMTPConfigurationService(base, smtpConfigurationPath, smtpIsConfiguredPath),
		Spaces:                         newSpaceService(base, spacesPath, spaceHomePath),
		Subscriptions:                  newSubscriptionService(base, subscriptionsPath),
		TagSets:                        newTagSetService(base, tagSetsPath, tagSetSortOrderPath),
		Tasks:                          newTaskService(base, tasksPath, taskTypesPath),
		TeamMembership:                 newTeamMembershipService(base, teamMembershipPath, teamMembershipPreviewTeamPath),
		Teams:                          newTeamService(base, teamsPath),
		Tenants:                        newTenantService(base, tenantsPath, tenantsMissingVariablesPath, tenantsStatusPath, tenantTagTestPath),
		TenantVariables:                newTenantVariableService(base, tenantVariablesPath),
		UpgradeConfiguration:           newUpgradeConfigurationService(base, upgradeConfigurationPath),
		UserOnboarding:                 newUserOnboardingService(base, userOnboardingPath),
		UserRoles:                      newUserRoleService(base, userRolesPath),
		Users:                          newUserService(base, usersPath, apiKeysPath, authenticateOctopusIDPath, currentUserPath, externalUserSearchPath, registerPath, signInPath, signOutPath, userAuthenticationPath, userIdentityMetadataPath),
		Variables:                      newVariableService(base, variablesPath, variableNamesPath, variablePreviewPath),
		WorkerPools:                    newWorkerPoolService(base, workerPoolsPath, workerPoolsDynamicWorkerTypesPath, workerPoolsSortOrderPath, workerPoolsSummaryPath, workerPoolsSupportedTypesPath),
		Workers:                        newWorkerService(base, workersPath, discoverWorkerPath, workerOperatingSystemsPath, workerShellsPath),
		WorkerToolsLatestImages:        newWorkerToolsLatestImageService(base, workerToolsLatestImagesPath),
	}, nil
}

// APIError is a generic structure for containing errors for API operations.
type APIError struct {
	Details         string   `json:"Details,omitempty"`
	ErrorMessage    string   `json:"ErrorMessage,omitempty"`
	Errors          []string `json:"Errors,omitempty"`
	FullException   string   `json:"FullException,omitempty"`
	HelpLink        string   `json:"HelpLink,omitempty"`
	HelpText        string   `json:"HelpText,omitempty"`
	ParsedHelpLinks []string `json:"ParsedHelpLinks,omitempty"`
}

// Error creates a predefined error for Octopus API responses.
func (e APIError) Error() string {
	return fmt.Sprintf("Octopus API error: %v %+v %v", e.ErrorMessage, e.Errors, e.FullException)
}

// APIErrorChecker is a generic error handler for the OctopusDeploy API.
func APIErrorChecker(urlPath string, resp *http.Response, wantedResponseCode int, slingError error, octopusDeployError *APIError) error {
	if octopusDeployError.Errors != nil {
		return fmt.Errorf("octopus deploy api returned an error on endpoint %s - %s", urlPath, octopusDeployError.Errors)
	}

	if slingError != nil {
		return fmt.Errorf("cannot get endpoint %s from server. failure from http client %v", urlPath, slingError)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrItemNotFound
	}

	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("bad request from endpoint %s. response from server %s", urlPath, resp.Status)
	}

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	if resp.StatusCode != wantedResponseCode {
		return octopusDeployError
	}

	return nil
}

// LoadNextPage checks if the next page should be loaded from the API. Returns
// the new path and a bool if the next page should be checked.
func LoadNextPage(pagedResults PagedResults) (string, bool) {
	if pagedResults.Links.PageNext != emptyString {
		return pagedResults.Links.PageNext, true
	}

	return emptyString, false
}

// Generic OctopusDeploy API Get Function.
func apiGet(sling *sling.Sling, inputStruct interface{}, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError(OperationAPIGet, ParameterSling)
	}

	getClient := sling.New()
	if getClient == nil {
		return nil, createClientInitializationError(OperationAPIGet)
	}

	getClient = getClient.Get(path)
	if getClient == nil {
		return nil, createClientInitializationError(OperationAPIGet)
	}

	getClient.Set("User-Agent", "go-octopusdeploy")

	octopusDeployError := new(APIError)
	resp, err := getClient.Receive(inputStruct, &octopusDeployError)
	if err != nil {
		return nil, err
	}

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)
	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return inputStruct, nil
}

// Generic OctopusDeploy API Add Function. Expects a 201 response.
func apiAdd(sling *sling.Sling, inputStruct interface{}, resource interface{}, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError(OperationAPIAdd, ParameterSling)
	}

	if isEmpty(path) {
		return nil, createInvalidParameterError(OperationAPIAdd, ParameterPath)
	}

	postClient := sling.New()
	if postClient == nil {
		return nil, createClientInitializationError(OperationAPIAdd)
	}

	postClient = postClient.Post(path)
	if postClient == nil {
		return nil, createClientInitializationError(OperationAPIAdd)
	}

	postClient.Set("User-Agent", "go-octopusdeploy")

	request := postClient.BodyJSON(inputStruct)
	if request == nil {
		return nil, createClientInitializationError(OperationAPIAdd)
	}

	octopusDeployError := new(APIError)
	resp, err := request.Receive(resource, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusCreated, err, octopusDeployError)
	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return resource, nil
}

// apiPost post to octopus and expect a 200 response code.
func apiPost(sling *sling.Sling, inputStruct interface{}, resource interface{}, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError(OperationAPIPost, ParameterSling)
	}

	if isEmpty(path) {
		return nil, createInvalidParameterError(OperationAPIPost, ParameterPath)
	}

	postClient := sling.New()
	if postClient == nil {
		return nil, createClientInitializationError(OperationAPIPost)
	}

	postClient = postClient.Post(path)
	if postClient == nil {
		return nil, createClientInitializationError(OperationAPIPost)
	}

	postClient.Set("User-Agent", "go-octopusdeploy")

	request := postClient.BodyJSON(inputStruct)
	if request == nil {
		return nil, createClientInitializationError(OperationAPIPost)
	}

	octopusDeployError := new(APIError)
	resp, err := request.Receive(resource, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)
	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return resource, nil
}

// Generic OctopusDeploy API Update Function.
func apiUpdate(sling *sling.Sling, inputStruct interface{}, resource interface{}, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError(OperationAPIUpdate, ParameterSling)
	}

	if isEmpty(path) {
		return nil, createInvalidParameterError(OperationAPIUpdate, ParameterPath)
	}

	putClient := sling.New()
	if putClient == nil {
		return nil, createClientInitializationError(OperationAPIUpdate)
	}

	putClient = putClient.Put(path)
	if putClient == nil {
		return nil, createClientInitializationError(OperationAPIUpdate)
	}

	putClient.Set("User-Agent", "go-octopusdeploy")

	request := putClient.BodyJSON(inputStruct)
	if request == nil {
		return nil, createClientInitializationError(OperationAPIUpdate)
	}

	octopusDeployError := new(APIError)
	resp, err := request.Receive(resource, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)
	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return resource, nil
}

// Generic OctopusDeploy API Delete Function.
func apiDelete(sling *sling.Sling, path string) error {
	if sling == nil {
		return createInvalidParameterError(OperationAPIDelete, ParameterSling)
	}

	if isEmpty(path) {
		return createInvalidParameterError(OperationAPIDelete, ParameterPath)
	}

	deleteClient := sling.New()
	if deleteClient == nil {
		return createClientInitializationError(OperationAPIDelete)
	}

	deleteClient = deleteClient.Delete(path)
	if deleteClient == nil {
		return createClientInitializationError(OperationAPIDelete)
	}

	deleteClient.Set("User-Agent", "go-octopusdeploy")

	octopusDeployError := new(APIError)
	resp, err := deleteClient.Receive(nil, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)
	if apiErrorCheck !=
		nil {
		return apiErrorCheck
	}

	return nil
}

// ErrItemNotFound is an OctopusDeploy error returned an item cannot be found.
var ErrItemNotFound = errors.New("cannot find the item")
