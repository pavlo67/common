package _demo

//func TestGeneratingEmptyList(t *testing.T) {
//
//	numberedFiles, err := filelib.NumberedFilesSequence(framesPath, framesRegexpStr, false)
//	require.NoError(t, err)
//	t.Logf("%d files found in (%s --- %s)", len(numberedFiles), framesPath, framesRegexpStr)
//
//	err = filelib.RemoveFile(listFilename)
//	require.NoError(t, err)
//
//	for i, numberedFile := range numberedFiles {
//		if i%10 == 0 {
//			fmt.Printf("%d/%d\n", i, len(numberedFiles))
//		}
//
//		frameMapping := position_log.Mapping{
//			Processing: video.Processing{
//				FrameInfo: video.FrameInfo{
//					N: numberedFile.I,
//				},
//			},
//			PointsRaw: []frame.PointRawGeo{{}},
//		}
//
//		err = serialization.SavePart(frameMapping, serialization.MarshalerJSON, listFilename)
//		require.NoError(t, err)
//	}
//
//}
