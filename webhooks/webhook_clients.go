package webhooks

// type testWebhookRequester[RequestType YokiTaskRunResult, ResponseType YokiTaskRunWebhookResponse] struct {
// }

// var currentWebhookClient WebhookRequester[YokiTaskRunResult, YokiTaskRunWebhookResponse]

// func (r *testWebhookRequester[in_T, out_T]) Request(url string, args *in_T) (*out_T, error) {
// 	var typed = (YokiTaskRunResult(*args))

// 	task, err := common.SelectYokiTask(typed.TaskId)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	time.Sleep(common.GetTestParams(task).WebhookCallDelay)

// 	req := &webhooks.HttpRequester[in_T, out_T]{}
// 	return req.Request(url, args)
// }

// func getWebhookClient() WebhookRequester[YokiTaskRunResult, YokiTaskRunWebhookResponse] {
// 	if currentWebhookClient == nil {
// 		if config.IsInTests() {
// 			currentWebhookClient = &testWebhookRequester[YokiTaskRunResult, YokiTaskRunWebhookResponse]{}
// 		} else {
// 			currentWebhookClient = &HttpRequester[YokiTaskRunResult, YokiTaskRunWebhookResponse]{}
// 		}

// 	}
// 	return currentWebhookClient
// }
