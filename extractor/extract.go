package extractor

import (
	sdk "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdk2 "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tf "github.com/hashicorp/terraform/helper/schema"
)

func ExtractProvider(provider interface{}, pi *ProviderInfo, outputPath string) {
	tfp, ok := provider.(*tf.Provider)
	if ok {
		extractor := new(Extractor)
		extractor.Generate(tfp, pi, outputPath)
	}
	sdkp, ok := provider.(*sdk.Provider)
	if ok {
		extractor := new(SdkExtractor)
		extractor.Generate(sdkp, pi, outputPath)
	}
	sdk2p, ok := provider.(*sdk2.Provider)
	if ok {
		extractor := new(Sdk2Extractor)
		extractor.Generate(sdk2p, pi, outputPath)
	}

}
