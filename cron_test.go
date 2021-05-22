package chrono

import (
	"testing"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

func TestCronExpression_NextTime(t *testing.T) {

	testCases := []struct {
		expression string
		time       string
		nextTimes  []string
	}{
		{
			"* * * * * *",
			"2021-05-31 23:59:56",
			[]string{
				"2021-05-31 23:59:57",
				"2021-05-31 23:59:58",
				"2021-05-31 23:59:59",
				"2021-06-01 00:00:00",
				"2021-06-01 00:00:01",
				"2021-06-01 00:00:02",
			},
		},
		{
			"17/3 * * * * *",
			"2021-03-16 15:04:16",
			[]string{
				"2021-03-16 15:04:17",
				"2021-03-16 15:04:20",
				"2021-03-16 15:04:23",
				"2021-03-16 15:04:26",
				"2021-03-16 15:04:29",
				"2021-03-16 15:04:32",
			},
		},
		{
			"19/3 * * * * *",
			"2021-03-16 15:04:19",
			[]string{
				"2021-03-16 15:04:22",
				"2021-03-16 15:04:25",
				"2021-03-16 15:04:28",
				"2021-03-16 15:04:31",
				"2021-03-16 15:04:34",
				"2021-03-16 15:04:37",
			},
		},
		{
			"8-19/3 * * * * *",
			"2021-03-16 15:04:23",
			[]string{
				"2021-03-16 15:05:08",
				"2021-03-16 15:05:11",
				"2021-03-16 15:05:14",
				"2021-03-16 15:05:17",
				"2021-03-16 15:06:08",
				"2021-03-16 15:06:11",
			},
		},
		{
			"8-24 * * * * *",
			"2021-03-16 15:04:23",
			[]string{
				"2021-03-16 15:04:24",
				"2021-03-16 15:05:08",
				"2021-03-16 15:05:09",
				"2021-03-16 15:05:10",
				"2021-03-16 15:05:11",
				"2021-03-16 15:05:12",
			},
		},
		{
			"0 * * * * *",
			"2021-05-21 13:41:37",
			[]string{
				"2021-05-21 13:42:00",
				"2021-05-21 13:43:00",
				"2021-05-21 13:44:00",
				"2021-05-21 13:45:00",
				"2021-05-21 13:46:00",
				"2021-05-21 13:47:00",
			},
		},
		{
			"7 * * * * *",
			"2021-05-22 13:12:56",
			[]string{
				"2021-05-22 13:13:07",
				"2021-05-22 13:14:07",
				"2021-05-22 13:15:07",
				"2021-05-22 13:16:07",
				"2021-05-22 13:17:07",
				"2021-05-22 13:18:07",
			},
		},
		{
			"0 0 * * * *",
			"2021-05-21 13:41:37",
			[]string{
				"2021-05-21 14:00:00",
				"2021-05-21 15:00:00",
				"2021-05-21 16:00:00",
				"2021-05-21 17:00:00",
				"2021-05-21 18:00:00",
				"2021-05-21 19:00:00",
			},
		},
		{
			"18 15 * * * *",
			"2021-05-21 19:12:56",
			[]string{
				"2021-05-21 19:15:18",
				"2021-05-21 20:15:18",
				"2021-05-21 21:15:18",
				"2021-05-21 22:15:18",
				"2021-05-21 23:15:18",
				"2021-05-22 00:15:18",
			},
		},
		{
			"18 15/5 * * * *",
			"2021-05-21 19:43:56",
			[]string{
				"2021-05-21 19:45:18",
				"2021-05-21 19:50:18",
				"2021-05-21 19:55:18",
				"2021-05-21 20:15:18",
				"2021-05-21 20:20:18",
				"2021-05-21 20:25:18",
			},
		},
		{
			"18 15-30/5 * * * *",
			"2021-05-21 19:43:56",
			[]string{
				"2021-05-21 20:15:18",
				"2021-05-21 20:20:18",
				"2021-05-21 20:25:18",
				"2021-05-21 20:30:18",
				"2021-05-21 21:15:18",
				"2021-05-21 21:20:18",
			},
		},
		{
			"18 40-45 * * * *",
			"2021-05-21 19:43:56",
			[]string{
				"2021-05-21 19:44:18",
				"2021-05-21 19:45:18",
				"2021-05-21 20:40:18",
				"2021-05-21 20:41:18",
				"2021-05-21 20:42:18",
				"2021-05-21 20:43:18",
			},
		},
		{
			"0 0 0 * * *",
			"2020-02-27 13:41:37",
			[]string{
				"2020-02-28 00:00:00",
				"2020-02-29 00:00:00",
				"2020-03-01 00:00:00",
				"2020-03-02 00:00:00",
				"2020-03-03 00:00:00",
				"2020-03-04 00:00:00",
			},
		},
		{
			"45 13 14 * * *",
			"2020-12-28 13:41:37",
			[]string{
				"2020-12-28 14:13:45",
				"2020-12-29 14:13:45",
				"2020-12-30 14:13:45",
				"2020-12-31 14:13:45",
				"2021-01-01 14:13:45",
				"2021-01-02 14:13:45",
			},
		},
		{
			"45 13 14/3 * * *",
			"2020-12-28 13:41:37",
			[]string{
				"2020-12-28 14:13:45",
				"2020-12-28 17:13:45",
				"2020-12-28 20:13:45",
				"2020-12-28 23:13:45",
				"2020-12-29 14:13:45",
				"2020-12-29 17:13:45",
			},
		},
		{
			"45 13 9-16/3 * * *",
			"2020-12-28 13:41:37",
			[]string{
				"2020-12-28 15:13:45",
				"2020-12-29 09:13:45",
				"2020-12-29 12:13:45",
				"2020-12-29 15:13:45",
				"2020-12-30 09:13:45",
				"2020-12-30 12:13:45",
			},
		},
		{
			"45 13 9-16 * * *",
			"2020-12-28 13:41:37",
			[]string{
				"2020-12-28 14:13:45",
				"2020-12-28 15:13:45",
				"2020-12-28 16:13:45",
				"2020-12-29 09:13:45",
				"2020-12-29 10:13:45",
				"2020-12-29 11:13:45",
			},
		},
		{
			"20 45 18 6 * *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-04-06 18:45:20",
				"2020-05-06 18:45:20",
				"2020-06-06 18:45:20",
				"2020-07-06 18:45:20",
				"2020-08-06 18:45:20",
				"2020-09-06 18:45:20",
			},
		},
		{
			"20 45 18 10-12 * *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-04-10 18:45:20",
				"2020-04-11 18:45:20",
				"2020-04-12 18:45:20",
				"2020-05-10 18:45:20",
				"2020-05-11 18:45:20",
				"2020-05-12 18:45:20",
			},
		},
		{
			"20 45 18 5-20/3 * *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-04-05 18:45:20",
				"2020-04-08 18:45:20",
				"2020-04-11 18:45:20",
				"2020-04-14 18:45:20",
				"2020-04-17 18:45:20",
				"2020-04-20 18:45:20",
			},
		},
		{
			"0 0 0 1 * *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-04-01 00:00:00",
				"2020-05-01 00:00:00",
				"2020-06-01 00:00:00",
				"2020-07-01 00:00:00",
				"2020-08-01 00:00:00",
				"2020-09-01 00:00:00",
			},
		},
		{
			"0 0 0 1 1 *",
			"2020-03-27 13:41:37",
			[]string{
				"2021-01-01 00:00:00",
				"2022-01-01 00:00:00",
				"2023-01-01 00:00:00",
				"2024-01-01 00:00:00",
				"2025-01-01 00:00:00",
				"2026-01-01 00:00:00",
			},
		},
		{
			"0 0 0 1 6 *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-06-01 00:00:00",
				"2021-06-01 00:00:00",
				"2022-06-01 00:00:00",
				"2023-06-01 00:00:00",
				"2024-06-01 00:00:00",
				"2025-06-01 00:00:00",
			},
		},
		{
			"0 0 0 1 3-12 *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-04-01 00:00:00",
				"2020-05-01 00:00:00",
				"2020-06-01 00:00:00",
				"2020-07-01 00:00:00",
				"2020-08-01 00:00:00",
				"2020-09-01 00:00:00",
			},
		},
		{
			"0 0 0 1 3-12/3 *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-06-01 00:00:00",
				"2020-09-01 00:00:00",
				"2020-12-01 00:00:00",
				"2021-03-01 00:00:00",
				"2021-06-01 00:00:00",
				"2021-09-01 00:00:00",
			},
		},
		{
			"0 0 0 1 SEP *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-09-01 00:00:00",
				"2021-09-01 00:00:00",
				"2022-09-01 00:00:00",
				"2023-09-01 00:00:00",
				"2024-09-01 00:00:00",
				"2025-09-01 00:00:00",
			},
		},
		{
			"0 0 0 1 AUG-OCT *",
			"2020-03-27 13:41:37",
			[]string{
				"2020-08-01 00:00:00",
				"2020-09-01 00:00:00",
				"2020-10-01 00:00:00",
				"2021-08-01 00:00:00",
				"2021-09-01 00:00:00",
				"2021-10-01 00:00:00",
			},
		},
		{
			"0 0 0 1 5 0",
			"2021-05-23 13:41:37",
			[]string{
				"2020-08-01 00:00:00",
				"2020-09-01 00:00:00",
				"2020-10-01 00:00:00",
				"2021-08-01 00:00:00",
				"2021-09-01 00:00:00",
				"2021-10-01 00:00:00",
			},
		},
	}

	for _, testCase := range testCases {
		exp, err := ParseCronExpression(testCase.expression)

		if err != nil {
			t.Errorf("could not parse cron expression : %s", err.Error())
			return
		}

		date, err := time.Parse(timeLayout, testCase.time)

		if err != nil {
			t.Errorf("could not parse time : %s", testCase.time)
			return
		}

		for _, nextTimeStr := range testCase.nextTimes {
			nextTime, err := time.Parse(timeLayout, nextTimeStr)

			if err != nil {
				t.Errorf("could not parse next time : %s", nextTimeStr)
				return
			}

			date = exp.NextTime(date)

			if nextTime.Format(timeLayout) != date.Format(timeLayout) {
				t.Errorf("got: %s expected: %s", date, nextTime)
			}
		}
	}

}
