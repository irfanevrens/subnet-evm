// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements ping tests, requires network-runner cluster.
package ping

//var _ = utils.DescribeLocal("[Ping]", func() {
//	ginkgo.It("can ping network-runner RPC server", func() {
//		runnerCli := runner.GetClient()
//		gomega.Expect(runnerCli).ShouldNot(gomega.BeNil())
//
//		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
//		_, err := runnerCli.Ping(ctx)
//		cancel()
//		gomega.Expect(err).Should(gomega.BeNil())
//	})
//})
//
