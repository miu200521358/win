package win

import (
	"syscall"
)

// IsWindowCenterObscured はウィンドウが最小化されているか、
// またはウィンドウの中心が他のアプリに隠れているかを判定します
func IsWindowCenterObscured(hwnd HWND) bool {
	// ウィンドウが最小化されているかをチェック
	if IsIconic(hwnd) {
		return true
	}

	// ウィンドウが表示されているかをチェック
	if !IsWindowVisible(hwnd) {
		return true
	}

	// ウィンドウの矩形を取得
	var rect RECT
	if !GetWindowRect(hwnd, &rect) {
		return false // 情報取得に失敗した場合はfalseを返す
	}

	// ウィンドウの中心点を計算
	centerX := (rect.Left + rect.Right) / 2
	centerY := (rect.Top + rect.Bottom) / 2

	// Z順（重なり順）でウィンドウを列挙し、対象の中心点を含む
	// より前面のウィンドウがあるかをチェックする
	var foundOverlappingWindow bool
	EnumWindows(syscall.NewCallback(func(hwndTop HWND, lParam uintptr) uintptr {
		// 自分自身は除外
		if hwndTop == hwnd {
			return 1 // TRUE相当
		}

		// 表示されていないウィンドウは無視
		if !IsWindowVisible(hwndTop) {
			return 1 // TRUE相当
		}

		// ウィンドウスタイルを取得し、子ウィンドウは無視
		style := GetWindowLong(hwndTop, GWL_STYLE)
		if (style & WS_CHILD) != 0 {
			return 1 // TRUE相当
		}

		// ウィンドウの矩形を取得
		var topRect RECT
		if GetWindowRect(hwndTop, &topRect) {
			// 対象の中心点がこのウィンドウの矩形内にあるかチェック
			if centerX >= topRect.Left && centerX <= topRect.Right &&
				centerY >= topRect.Top && centerY <= topRect.Bottom {
				// Z順を比較する前に位置関係を直接チェックする方法を使用
				if IsWindowInFront(hwndTop, hwnd) {
					foundOverlappingWindow = true
					return 0 // FALSE相当、列挙を中止
				}
			}
		}

		return 1 // TRUE相当
	}), 0)

	return foundOverlappingWindow
}

// IsWindowInFront は window1 が window2 より前面にあるかをより効率的にチェックします
func IsWindowInFront(window1, window2 HWND) bool {
	// GetNextWindow を使って前面から順にウィンドウを取得
	hwndCur := GetWindow(GetDesktopWindow(), GW_CHILD)

	// window1 と window2 のいずれかが見つかった時点で、どちらが先かで判定
	for hwndCur != 0 {
		if hwndCur == window1 {
			return true // window1 が先に見つかった = 前面にある
		}
		if hwndCur == window2 {
			return false // window2 が先に見つかった = 後面にある
		}
		hwndCur = GetWindow(hwndCur, GW_HWNDNEXT)
	}

	return false // どちらも見つからなかった場合
}

// GetWindowZOrder はウィンドウのZ順（重なり順）を返します
// 値が小さいほど前面にあることを表します
func GetWindowZOrder(hwnd HWND) int {
	zOrder := 0
	hwndTop := GetDesktopWindow()
	hwndCur := GetWindow(hwndTop, GW_CHILD)

	for hwndCur != 0 {
		if hwndCur == hwnd {
			return zOrder
		}
		zOrder++
		hwndCur = GetWindow(hwndCur, GW_HWNDNEXT)
	}

	return 99999 // 見つからなかった場合は大きな値を返す
}
