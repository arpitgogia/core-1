package bhattacharya

import (
	"coralreefci/models/issues"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// DOC: helper function for cleaning up calculated results.
func Round(input float64) float64 {
	rounded := math.Floor((input*10000.0)+0.5) / 10000.0
	return rounded
}

// DOC: helper function to translate float64 into string.
func ToString(number float64) string {
	return strconv.FormatFloat(number, 'f', 4, 64)
}

// DOC: argument inputs into FoldImplementation:
//      trainIssues - this is the fold-defined length to train on (e.g. 10%)
//      testIssues - this is the fold-defined length to test on (e.g. 90%)
func (m *Model) FoldImplementation(trainIssues, testIssues []issues.Issue) (float64, matrix) {
	testCount, correct := len(testIssues), 0
	expected, predicted := []issues.Issue{}, []issues.Issue{}
	expected = append(expected, testIssues...)
	predicted = append(predicted, testIssues...)
	m.Learn(trainIssues)

	for i := 0; i < len(testIssues); i++ {
		assignees := m.Predict(testIssues[i]) // NOTE: assignees is working appropriately
		for j := 0; j < len(assignees); j++ {
			// NOTE: This is not tested logging functionality
			// model.Logger.Log(testIssues[j].Url)
			// model.Logger.Log(testIssues[j].Assignee)
			if assignees[j] == testIssues[i].Assignee {
				correct++
				predicted[i].Assignee = assignees[j]
				break
			} else {
				predicted[i].Assignee = assignees[0]
			}
		}
	}
	mat, err := BuildMatrix(expected, predicted)
	if err != nil {
		fmt.Println(err)
	}
	return float64(correct) / float64(testCount), mat
}

func (m *Model) JohnFold(issues []issues.Issue) (string, error) {
	issueCount := len(issues)
	if issueCount < 10 {
		return "", errors.New("LESS THAN 10 ISSUES SUBMITTED - JOHN FOLD")
	}

	finalScore := 0.00
	for i := 0.10; i < 0.90; i += 0.10 { // TODO: double check the logic / math here
		trainCount := int(Round(i * float64(issueCount)))

		// TODO: add in logging here for the output matrix on each loop run
		fmt.Println("\n\nFOLD TRAIN SIZE ", i)
		score, _ := m.FoldImplementation(issues[:trainCount], issues[trainCount:])

		finalScore += score
	}
	return ToString(Round(finalScore / 9.00)), nil
}

func (m *Model) TwoFold(issueList []issues.Issue) (string, error) {
	issueCount := len(issueList)
	if issueCount < 10 {
		return "", errors.New("LESS THAN 10 ISSUES SUBMITTED - TWO FOLD")
	}
	trainEndPos := int(0.50 * float64(issueCount))
	trainIssues := []issues.Issue{} // TODO: put into one line (see #84)
	testIssues := []issues.Issue{}  // TODO: put into one line (see #83)
	trainIssues = append(trainIssues, issueList[0:trainEndPos]...)
	testIssues = append(testIssues, issueList[trainEndPos+1:]...)

	firstScore, _ := m.FoldImplementation(trainIssues, testIssues)
	secondScore, _ := m.FoldImplementation(testIssues, trainIssues)

	score := firstScore + secondScore

	return ToString(Round(score / 2.00)), nil
}

func (m *Model) TenFold(issueList []issues.Issue) (string, error) {
	issueCount := len(issueList)
	if issueCount < 10 {
		return "", errors.New("LESS THAN 10 ISSUES SUBMITTED - TEN FOLD")
	}

	finalScore := 0.00
	start := 0
	for i := 0.10; i <= 1.00; i += 0.10 {
		end := int(Round(i * float64(issueCount)))

		segment := issueList[start:end]
		remainder := []issues.Issue{}
		remainder = append(issueList[:start], issueList[end:]...)

		score, _ := m.FoldImplementation(segment, remainder) // TODO: specific logs for matrices

		finalScore += score
		start = end
	}
	return ToString(Round(finalScore / 10.00)), nil
}
