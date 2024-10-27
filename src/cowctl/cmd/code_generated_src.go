package cmd

import complete "github.com/chriswalz/complete/v3"

var rootSuggestionTree = &complete.CompTree{Desc: "cowctl", Sub: map[string]*complete.CompTree{
	"init": {Desc: "initialize rule/task", Sub: map[string]*complete.CompTree{
		"rule": {Desc: "initialize rule", Flags: map[string]*complete.CompTree{
			"path":        {Desc: `path to the folder to initialize the rule`},
			"name":        {Desc: `name of the rule`},
			"config-path": {Desc: `path to the configuration file(yaml/json)`},
			"rules-path":  {Desc: `path to the rules folder`},
			"tasks-path":  {Desc: `path to the tasks folder`},
			"exec-path":   {Desc: `path to the maintain the executions`},
			"catalog":     {Desc: `default "localcatalog/rules" use "globalcatalog" to init rule in "globalcatalog/rules"`},
		}},
		"task": {Desc: "initialize task", Flags: map[string]*complete.CompTree{
			"path":        {Desc: `path to the folder to initialize the task`},
			"name":        {Desc: `name of the task`},
			"language":    {Desc: `programming language - go/python(default-go)`},
			"config-path": {Desc: `path to the configuration file(yaml/json)`},
			"rules-path":  {Desc: `path to the rules folder`},
			"tasks-path":  {Desc: `path to the tasks folder`},
			"exec-path":   {Desc: `path to the maintain the executions`},
			"catalog":     {Desc: `default "localcatalog/tasks" use "globalcatalog" to init task in "globalcatalog/tasks"`},
		}},
		"credential": {Desc: "initialize credential", Flags: map[string]*complete.CompTree{
			"path":        {Desc: `path to the folder to initialize the credential`},
			"name":        {Desc: `name of the credential`},
			"config-path": {Desc: `path to the credential file(yaml/json)`},
		}},
		"application": {Desc: "initialize application", Flags: map[string]*complete.CompTree{
			"path":        {Desc: `path to the folder to initialize the application`},
			"name":        {Desc: `name of the application`},
			"config-path": {Desc: `path to the application file(yaml/json)`},
		}},
	},
	},
	"exec": {Desc: "execute the rules", Sub: map[string]*complete.CompTree{
		"rule": {Desc: "execute the rule", Flags: map[string]*complete.CompTree{
			"verbose":                  {Desc: `Display the rule outputs in the console`},
			"config-path":              {Desc: `path to the configuration file(yaml/json)`},
			"rules-path":               {Desc: `path to the rules folder`},
			"tasks-path":               {Desc: `path to the tasks folder`},
			"exec-path":                {Desc: `path to the maintain the executions`},
			"rule-name":                {Desc: `name of the rule to be execute`},
			"downloads-path":           {Desc: `path of the downloads folder`},
			"catalog":                  {Desc: `default "localcatalog/rules" use "globalcatalog" search only "globalcatalog/rules" for rule`},
			"preserve-execution-setup": {Desc: `execution set up will be retained for unit test`},
		}},
		"rulegroup": {Desc: "execute the rulegroup", Flags: map[string]*complete.CompTree{
			"verbose":                  {Desc: `Display the rule outputs in the console`},
			"config-path":              {Desc: `path to the configuration file(yaml/json)`},
			"rules-path":               {Desc: `path to the rules folder`},
			"tasks-path":               {Desc: `path to the tasks folder`},
			"exec-path":                {Desc: `path to the maintain the executions`},
			"rule-group-path":          {Desc: `path to the rule group`},
			"rule-group-name":          {Desc: `name of the rule group to be execute`},
			"downloads-path":           {Desc: `path of the downloads folder`},
			"catalog":                  {Desc: `default "localcatalog/rulegroups" use "globalcatalog" search only "globalcatalog/rulegroups" for rulegroup`},
			"preserve-execution-setup": {Desc: `execution set up will be retained for unit test`},
		}},
	},
	},
	"export": {Desc: "export the rule/rulegroup", Sub: map[string]*complete.CompTree{
		"rule": {Desc: "export the rule", Flags: map[string]*complete.CompTree{
			"config-path":      {Desc: `path to the configuration file(yaml/json)`},
			"rules-path":       {Desc: `path to the rules folder`},
			"tasks-path":       {Desc: `path to the tasks folder`},
			"exec-path":        {Desc: `path to the maintain the executions`},
			"rule-name":        {Desc: `name of the rule to be execute`},
			"downloads-path":   {Desc: `path of the downloads folder`},
			"export-file-type": {Desc: `export file type(zip/tar). default tar`},
			"catalog":          {Desc: `default "localcatalog/rules" use "globalcatalog" search only "globalcatalog/rules" for rule`},
		}},
		"rulegroup": {Desc: "export the rulegroup", Flags: map[string]*complete.CompTree{
			"config-path":      {Desc: `path to the configuration file(yaml/json)`},
			"rules-path":       {Desc: `path to the rules folder`},
			"tasks-path":       {Desc: `path to the tasks folder`},
			"exec-path":        {Desc: `path to the maintain the executions`},
			"rule-group-path":  {Desc: `path to the rule group`},
			"rule-group-name":  {Desc: `name of the rule group to be execute`},
			"downloads-path":   {Desc: `path of the downloads folder`},
			"export-file-type": {Desc: `export file type(zip/tar). default tar`},
		}},
	},
	},
	"publish": {Desc: "publish the rule/appconfig/credconfig", Sub: map[string]*complete.CompTree{
		"rule": {Desc: "publish the rule", Flags: map[string]*complete.CompTree{
			"config-path":         {Desc: `path to the configuration file(yaml/json)`},
			"rules-path":          {Desc: `path to the rules folder`},
			"tasks-path":          {Desc: `path to the tasks folder`},
			"exec-path":           {Desc: `path to the maintain the executions`},
			"rule-name":           {Desc: `name of the rule to be published`},
			"downloads-path":      {Desc: `path of the downloads folder`},
			"export-file-type":    {Desc: `export file type(zip/tar). default tar`},
			"publish-name":        {Desc: `this value used as rule name while creating in compliance cow`},
			"publish-description": {Desc: `this value used as rule description while creating in compliance cow`},
			"client-id":           {Desc: `Id which is generated by the user`},
			"client-secret":       {Desc: `Secret key which is generated by user`},
			"sub-domain":          {Desc: `where to publish? dev/partner. default partner `},
			"user-domain":         {Desc: `Where should the content be published if the user is connected to multiple domains? Default: Primary domain.`},
			"catalog":             {Desc: `default "localcatalog/rules" use "globalcatalog" search only "globalcatalog/rules" for rule`},
		}},
		"application": {Desc: "publish application", Flags: map[string]*complete.CompTree{
			"name":          {Desc: `name of the application`},
			"version":       {Desc: `version of the application`},
			"config-path":   {Desc: `path to the configuration file(yaml/json)`},
			"client-id":     {Desc: `Id which is generated by the user`},
			"client-secret": {Desc: `Secret key which is generated by user`},
			"user-domain":   {Desc: `Where should the content be published if the user is connected to multiple domains? Default: Primary domain.`},
		}},
	},
	},
	"create": {Desc: "create appconfig/credconfig with yaml file", Sub: map[string]*complete.CompTree{
		"credential": {Desc: "create credential", Flags: map[string]*complete.CompTree{
			"config-path": {Desc: `path to the configuration file(yaml/json)`},
			"file-name":   {Desc: `name of the yaml file`},
		}},
		"application": {Desc: "create application", Flags: map[string]*complete.CompTree{
			"config-path": {Desc: `path to the configuration file(yaml/json)`},
			"file-name":   {Desc: `name of the yaml file`},
		}},
	}},
	"exit": {Desc: "exit the cli terminal", Sub: map[string]*complete.CompTree{}},
}}
