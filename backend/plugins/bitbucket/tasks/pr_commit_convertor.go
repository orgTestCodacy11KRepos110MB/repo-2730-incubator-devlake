/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"reflect"

	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var ConvertPrCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertPullRequestCommits",
	EntryPoint:       ConvertPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_pull_request_commits into  domain layer table pull_request_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func ConvertPullRequestCommits(taskCtx plugin.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	repoId := data.Repo.BitbucketId

	pullIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketPullRequest{})

	cursor, err := db.Cursor(
		dal.From(&bitbucketModels.BitbucketPrCommit{}),
		dal.Join(`left join _tool_bitbucket_pull_requests on _tool_bitbucket_pull_requests.bitbucket_id = _tool_bitbucket_pull_request_commits.pull_request_id`),
		dal.Where("_tool_bitbucket_pull_requests.repo_id = ? and _tool_bitbucket_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
		dal.Orderby("pull_request_id ASC"),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketPrCommit{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_PULL_REQUEST_COMMITS_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketPullRequestCommit := inputRow.(*bitbucketModels.BitbucketPrCommit)
			domainPrCommit := &code.PullRequestCommit{
				CommitSha:     bitbucketPullRequestCommit.CommitSha,
				PullRequestId: pullIdGen.Generate(data.Options.ConnectionId, bitbucketPullRequestCommit.PullRequestId),
			}
			return []interface{}{
				domainPrCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
