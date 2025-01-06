package _demo

//func TestCreateRulers(t *testing.T) {
//
//	var framesMapping []position_log.Mapping
//	err := serialization.ReadAllPartsJSON(listFilename, &framesMapping)
//	require.NoError(t, err)
//
//	numberedFiles, err := filelib.NumberedFilesSequence(framesPath, framesRegexpStr, true)
//	require.NoError(t, err)
//	t.Logf("%d files found in (%s --- %s)", len(numberedFiles), framesPath, framesRegexpStr)
//
//	for _, numberedFile := range numberedFiles {
//		if numberedFile.I < 230 {
//			continue
//		}
//
//		fmt.Printf("%d\n", numberedFile.I)
//
//		if CheckIfFileIsSelected(framesMapping, numberedFile.I) >= 0 {
//			img, err := imagelib.ReadImage(numberedFile.Path)
//			require.NoError(t, err)
//			require.NotNil(t, img)
//
//			imgRGB, err := imagelib.ImageToRGBA(img)
//			require.NoError(t, err)
//			require.NotNil(t, imgRGB)
//
//			gridMarking := imagelib.GridMarking{
//				DividerX: dividerX,
//				DividerY: dividerY,
//			}
//
//			imagelib.MarkGrid(imgRGB, gridMarking, colornames.Red, colornames.White, 20)
//
//			rulersFile := filepath.Join(rulersPath, fmt.Sprintf(filenameTpl, numberedFile.I))
//			err = imagelib.SavePNG(imgRGB, rulersFile)
//			require.NoError(t, err)
//
//		} else {
//			err = filelib.RemoveFile(numberedFile.Path)
//			require.NoError(t, err)
//
//		}
//
//	}
//}
