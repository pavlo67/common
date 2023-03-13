package numlib

import (
	"fmt"
	"log"
	"sort"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

type ParameterInterval struct {
	Name    string
	Min     float64
	Max     float64
	Step    float64
	Divider int
}

type Parameter struct {
	Value       float64
	CurrentStep float64
}

type Result struct {
	Parameters   []Parameter
	Criterion    float64
	Restrictions []float64
	Details      interface{}
}

type Model interface {
	CheckParams([]ParameterInterval) ([]ParameterInterval, error)
	Count([]Parameter, common.Map) (Result, error)

	// SetComparativeResult(parameters []Parameter) (Result, error)
}

const parameterStepRatioDefault = 0.4

const onGradientSearch = "on GradientSearch()"

func GradientSearch(model Model, options common.Map, parameterIntervals []ParameterInterval, maxResults int, parameterStepRatio float64) ([]Result, error) {
	if model == nil {
		return nil, errors.New(onGradientSearch + ": model is nil")
	}

	parameterIntervals, err := model.CheckParams(parameterIntervals)
	if err != nil {
		return nil, errors.Wrap(err, onGradientSearch)
	}

	if len(parameterIntervals) < 1 {
		return nil, errors.New(onGradientSearch + ": parameterIntervals is empty")
	}

	if parameterStepRatio <= 0 || parameterStepRatio >= 1 {
		log.Printf("parameterStepRatio is changed from %f to default = %f", parameterStepRatio, parameterStepRatioDefault)
		parameterStepRatio = parameterStepRatioDefault
	}

	log.Printf("parameter intervals for gradient search: %#v", parameterIntervals)

	initialParametersSet := [][]Parameter{{}}

	for _, parameterInterval := range parameterIntervals {
		if parameterInterval.Max < parameterInterval.Min {
			parameterInterval.Min, parameterInterval.Max = parameterInterval.Max, parameterInterval.Min
		}

		if parameterInterval.Divider <= 1 {
			parameterInterval.Divider = 0
		} else {
			parameterInterval.Step = (parameterInterval.Max - parameterInterval.Min) / float64(parameterInterval.Divider)
		}

		var initialParametersSetNew [][]Parameter

		var value float64

		for n := 0; n <= parameterInterval.Divider; n++ {
			if parameterInterval.Divider == 0 {
				value = 0.5 * (parameterInterval.Max + parameterInterval.Min)
			} else if n == parameterInterval.Divider {
				value = parameterInterval.Max
			} else {
				value = parameterInterval.Min + float64(n)*parameterInterval.Step
			}

			for i := 0; i < len(initialParametersSet); i++ {
				initialParametersSetNew = append(initialParametersSetNew, append(initialParametersSet[i], Parameter{
					Value:       value,
					CurrentStep: parameterInterval.Step,
				}))
			}
		}

		initialParametersSet = initialParametersSetNew
	}

	var results []Result

	for i, initialParameters := range initialParametersSet {
		if i%1000 == 1 {
			fmt.Printf("PROCESSED %d BRANCHES OF TOTAL %d\n", i, len(initialParametersSet))
		}

		initialResult, err := model.Count(initialParameters, options)
		if err != nil {
			return results, errors.Wrap(err, onGradientSearch)
		}

		// TODO!!! be careful if model.Count() changes some parameters
		initialResult.Parameters = initialParameters

		result, err := GradientRecursion(model, initialResult, parameterStepRatio)
		if err != nil {
			return results, errors.Wrap(err, onGradientSearch)
		}

		results = append(results, result)

		// here is an error ???
		//// TODO: check limitations??? it should be done in model.Count()???
		//if maxResults <= 0 || len(results) < maxResults {
		//	results = append(results, result)
		//} else if results[len(results)-1].Criterion < result.Criterion {
		//	results[len(results)-1] = result
		//}

		sort.Slice(results, func(i, j int) bool { return results[i].Criterion < results[j].Criterion })
	}

	return results, nil
}

func GradientRecursion(model Model, result Result, parameterStepRatio float64) (Result, error) {
	// TODO!!!
	return result, nil
}

func VariateParameters(parameters []Parameter) [][]Parameter {
	return nil
}
