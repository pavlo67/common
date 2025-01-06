package _demo

import (
	"github.com/pavlo67/imagelib/sources"
	"path/filepath"
)

const rootPath = "1.2/"

// const framesPath = rootPath + sources.FramesSubPath
const framesRegexpStr = `^(\d{4})\.png$`
const rulersPath = rootPath + "_rulers/"
const filenameTpl = "%04d.png"

// var listFilename = filepath.Join(rootPath, sources.FramesJSONLog)
var listInterpolatedFilename = filepath.Join(rootPath, sources.FramesAllDescriptionsFilename)

const dividerX = 30
const dividerY = 20

//func TestCalculateFramePosition(t *testing.T) {
//	N := 285
//
//	var framesMapping []position_log.Mapping
//	err := serialization.ReadAllPartsJSON(listFilename, &framesMapping)
//	require.NoError(t, err)
//
//	if frameI := CheckIfFileIsSelected(framesMapping, N); frameI >= 0 {
//		frameFile := filepath.Join(framesPath, fmt.Sprintf(filenameTpl, N))
//
//		img, err := imagelib.ReadImage(frameFile)
//		require.NoError(t, err)
//		require.NotNil(t, img)
//
//		rect := img.Bounds()
//
//		xWidth, yHeight := rect.Dx(), rect.Dy()
//
//		unitX, unitY := float64(xWidth)/float64(dividerX), float64(yHeight)/float64(dividerY)
//
//		ptsRaw := [3]frame.PointRawGeo{
//			{
//				Grid: plane.Point2{19, 13.9},
//				Geo:  geolib.Point{48.550459, 38.136724},
//			},
//			{
//				Grid: plane.Point2{21.9, 18.4},
//				Geo:  geolib.Point{48.550966, 38.137454},
//			},
//			{
//				Grid: plane.Point2{0.7, 14.8},
//				Geo:  geolib.Point{48.547603, 38.137708},
//			},
//		}
//
//		geoPoint, rotation, dpmOriginal, err := frame.CalculateWithGeoPoints(xWidth, yHeight, ptsRaw, unitX, unitY)
//		require.NoError(t, err)
//		require.NotNil(t, geoPoint)
//		require.Truef(t, dpmOriginal > 0 && !math.IsInf(dpmOriginal, 1), "dpmOriginal: %f", dpmOriginal)
//
//		t.Logf("geoPoint: %v, dpm original: %f", geoPoint, dpmOriginal)
//
//		SaveFrame(t, img, dpmOriginal, *geoPoint, rotation, ptsRaw, framesMapping, frameI, N)
//	}
//}
//
//func TestInterpolateFramePositions(t *testing.T) {
//	err := filelib.RemoveFile(listInterpolatedFilename)
//	require.NoError(t, err)
//
//	var framesMapping []position_log.Mapping
//	err = serialization.ReadAllPartsJSON(listFilename, &framesMapping)
//	require.NoError(t, err)
//	require.Truef(t, len(framesMapping) >= 3, "framesMapping is too short to be interpolated: %d", len(framesMapping))
//
//	framesMappingInterpolated := []position_log.Mapping{framesMapping[0]}
//
//	for i, frameMapping := range framesMapping[:len(framesMapping)-1] {
//		minN, maxN := frameMapping.N, framesMapping[i+1].N
//		require.Truef(t, minN < maxN, "framesMapping is wrong sorted at %d (%+v) / %d (%+v)", i, frameMapping, i+1, framesMapping[i+1])
//
//		if framesMapping[i+1].Position == nil {
//			break
//		}
//
//		interval := framesMapping[i+1].Time.Sub(frameMapping.Time) / time.Duration(maxN-minN)
//
//		for n := minN + 1; n < maxN; n++ {
//			x, x0, x1 := float64(n), float64(minN), float64(maxN)
//
//			beaFloat, _ := mathlib.InterpolateByTwoPoints(x, [2][2]float64{{x0, float64(frameMapping.Position.Bearing)}, {x1, float64(framesMapping[i+1].Position.Bearing)}})
//			altitude, _ := mathlib.InterpolateByTwoPoints(x, [2][2]float64{{x0, frameMapping.Position.Altitude}, {x1, framesMapping[i+1].Position.Altitude}})
//			latFloat, _ := mathlib.InterpolateByTwoPoints(x, [2][2]float64{{x0, float64(frameMapping.Point.Lat)}, {x1, float64(framesMapping[i+1].Point.Lat)}})
//			lonFloat, _ := mathlib.InterpolateByTwoPoints(x, [2][2]float64{{x0, float64(frameMapping.Point.Lon)}, {x1, float64(framesMapping[i+1].Point.Lon)}})
//			dpm, _ := mathlib.InterpolateByTwoPoints(x, [2][2]float64{{x0, frameMapping.Settings.DPM}, {x1, framesMapping[i+1].Settings.DPM}})
//
//			frameMappingInterpolated := position_log.Mapping{
//				Processing: video.Processing{
//					FrameInfo: video.FrameInfo{
//						N: n,
//						Position: &positions.Item{
//							Bearing:  geolib.Bearing(beaFloat),
//							Altitude: altitude,
//							//Pitch:    0,
//							//Roll:     0,
//						},
//						Point: &geolib.Point{
//							Lat: geolib.Degrees(latFloat),
//							Lon: geolib.Degrees(lonFloat),
//						},
//						Settings: imagelib.Settings{
//							DPM: dpm,
//						},
//						Time: frameMapping.Time.Add(time.Duration(n-minN) * interval),
//					},
//				},
//			}
//			framesMappingInterpolated = append(framesMappingInterpolated, frameMappingInterpolated)
//		}
//
//		framesMappingInterpolated = append(framesMappingInterpolated, framesMapping[i+1])
//	}
//
//	err = serialization.SaveAllPartsJSON(framesMappingInterpolated, listInterpolatedFilename)
//	require.NoError(t, err)
//}
//
//func TestRecalculateFramePositions(t *testing.T) {
//	var framesMapping []position_log.Mapping
//	err := serialization.ReadAllPartsJSON(listFilename, &framesMapping)
//	require.NoError(t, err)
//
//	for frameI, frameMapping := range framesMapping {
//		t.Log(frameI)
//
//		if len(frameMapping.PointsRaw) != 3 {
//			continue
//		}
//
//		N := frameMapping.N
//		frameFile := filepath.Join(framesPath, fmt.Sprintf(filenameTpl, N))
//
//		img, err := imagelib.ReadImage(frameFile)
//		require.NoError(t, err)
//		require.NotNil(t, img)
//
//		rect := img.Bounds()
//
//		xWidth, yHeight := rect.Dx(), rect.Dy()
//
//		unitX, unitY := float64(xWidth)/float64(dividerX), float64(yHeight)/float64(dividerY)
//
//		ptsRaw := [3]frame.PointRawGeo{frameMapping.PointsRaw[0], frameMapping.PointsRaw[1], frameMapping.PointsRaw[2]}
//
//		geoPoint, rotation, dpmOriginal, err := frame.CalculateWithGeoPoints(xWidth, yHeight, ptsRaw, unitX, unitY)
//		require.NoError(t, err)
//		require.NotNil(t, geoPoint)
//		require.Truef(t, dpmOriginal > 0 && !math.IsInf(dpmOriginal, 1), "dpmOriginal: %f", dpmOriginal)
//
//		// t.Logf("geoPoint: %v, dpm original: %f", geoPoint, dpmOriginal)
//
//		SaveFrame(t, img, dpmOriginal, *geoPoint, rotation, ptsRaw, framesMapping, frameI, N)
//
//	}
//
//}
//
//func SaveFrame(t *testing.T, img image.Image, dpmOriginal float64, geoPoint geolib.Point, rotation plane.Rotation, ptsRaw [3]frame.PointRawGeo,
//	framesMapping []position_log.Mapping, frameI, N int) {
//	l := logger_test.New(t, strconv.FormatInt(time.Now().Unix(), 10), "", true, nil)
//	tilesOp := StartTilesOp(t, l)
//	sourceKey := tiles_arcgis.SourceKey
//
//	dpmRequired := 2.
//	zoom, dpm, err := tiles2.CalculateZoom(tilesOp, sourceKey, geoPoint.Lat, dpmRequired)
//	require.NoError(t, err)
//
//	rect := img.Bounds()
//
//	xWidth, yHeight := rect.Dx(), rect.Dy()
//
//	marginRatio := 1.1
//	sideX := float64(xWidth) * marginRatio / dpm
//	sideY := float64(yHeight) * marginRatio / dpm
//
//	l.Infof("dpm: %f,  sideX: %f, sideY: %f", dpm, sideX, sideY)
//
//	imgRGBTiled, dpm, err := tiles2.PositionImage(tilesOp, sourceKey, geoPoint, rotation, sideX, sideY, zoom, nil)
//	require.NoError(t, err)
//	require.Truef(t, dpm > 0 && !math.IsInf(dpm, 1), "dpmOriginal: %f", dpm)
//	require.NotNil(t, imgRGBTiled)
//
//	filenameFromTiles := filepath.Join(framesPath, fmt.Sprintf("%04d_from_tiles.png", N))
//	err = imagelib.SavePNG(imgRGBTiled, filenameFromTiles)
//	require.NoError(t, err)
//
//	pointsGrid := []plane.Point2{ptsRaw[0].Grid, ptsRaw[1].Grid, ptsRaw[2].Grid}
//
//	gridMarking := imagelib.GridMarking{
//		DividerX:   dividerX,
//		DividerY:   dividerY,
//		GridPoints: pointsGrid,
//	}
//
//	imgRGB, err := imagelib.ImageToRGBA(img)
//	require.NoError(t, err)
//	require.NotNil(t, imgRGB)
//
//	imagelib.MarkGrid(imgRGB, gridMarking, colornames.Red, colornames.Blue, 20)
//	filenameRuler := filepath.Join(rulersPath, fmt.Sprintf(filenameTpl, N))
//	err = imagelib.SavePNG(imgRGB, filenameRuler)
//	require.NoError(t, err)
//
//	framesMapping[frameI].PointsRaw = ptsRaw[:]
//	framesMapping[frameI].Position = &positions.Item{
//		Bearing: geolib.PlaneBearingFromRotation(rotation),
//	}
//	framesMapping[frameI].Settings = imagelib.Settings{DPM: dpmOriginal}
//	framesMapping[frameI].Point = &geoPoint
//
//	err = os.Rename(listFilename, listFilename+"."+time.Now().Format(time.RFC3339)+".bak")
//	require.NoError(t, err)
//
//	err = serialization.SaveAllPartsJSON(framesMapping, listFilename)
//	require.NoError(t, err)
//
//}
//
//func CheckIfFileIsSelected(framesMapping []position_log.Mapping, n int) int {
//	for i, frameMapping := range framesMapping {
//		if frameMapping.N == n {
//			return i
//		}
//	}
//
//	return -1
//}
//
//func StartTilesOp(t *testing.T, l logger.Operator) tiles2.Operator {
//
//	env, err := config.Get("../_env/local.yaml", serialization.MarshalerYAML)
//	require.NoError(t, err)
//	require.NotNil(t, env)
//
//	components := []starter.Component{
//		// common sensors ---------------------------
//		{tiles_files.Starter(), nil},
//		{tiles_arcgis.Starter(), nil},
//		{tiles_manager.Starter(), common.Map{"ext_interface_key": tiles_arcgis.InterfaceKey}},
//	}
//
//	label := "TILES/TEST BUILD"
//	joinerOp, err := starter.Run(components, env, label, l)
//	require.NoError(t, err)
//
//	defer joinerOp.CloseAll()
//
//	tilesOp, _ := joinerOp.Interface(tiles2.InterfaceKey).(tiles2.Operator)
//	require.NotNil(t, tilesOp)
//
//	return tilesOp
//}
