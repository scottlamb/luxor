// Package protocol defines the abstract RPC interface for the FX Luminaire
// Luxor ZD landscape lighting system's wi-fi module.
//
// Each Controller method is of the form:
//
// 	Method(ctx context.Context, request *MethodRequest) (response *MethodResponse, err error)
//
// and each behaves similarly: it makes the given RPC, using
// ctx for cancellation and deadline. If the RPC fails at either the
// protocol layer or via response.Status != 0, err will be non-nil.
//
// Important concepts:
//
// Group numbers: lights are semi-permanently assigned to uint8 "group numbers" (in the full
// [0, 256) range) in one of two ways:
//
// (1) Physically plugging the light into a special port on the controller.
//
// (2) The "AssignLight" command given a light's serial number. The serial
// number can be learned by attaching a special wireless dongle to a phone's
// audio port and listening for serial numbers via a separate protocol. The
// author of this package does not own such a dongle, so this is untested.
//
// There are zero or more lights per group number. Group numbers can be used
// to turn lights to a given uint8 intensity. Useful values are [0, 100];
// higher values may be set but do not produce brighter light than 100.
//
// Groups: Each group number has zero or one "group". Groups have distinct
// names and are displayed in the app's UI in a user-controllable order.
// Groups are persisted across controller restarts and can be changed through
// this protocol.
//
// Themes: A "theme" is a stored list of (group, intensity) tuples used to set
// groups to the given intensities or 0%. Apparently if a list contains a
// group more than once, the last such tuple wins.
//
// Themes have a name (up to 19 bytes) and a uint8 index in [0, 26), which
// corresponds to ['A', 'Z'] in the controller and app UI. It is apparently
// possible to have two themes with the same name and/or index, but following
// commands will be ambiguous, so this is not a desirable state.
//
// Themes are "on" or "off", matching the last IlluminateTheme (or
// ExtinguishAll) request. Because it is possible to change the theme's
// group's intensities in other ways, the lights do not always match the
// theme's state.
package protocol

import (
	"context"
	"errors"
	"fmt"
)

type Controller interface {
	// AssignLight assigns a light to a group number. Must know its serial
	// number. The app learns serial numbers via a separate protocol that
	// uses a special dongle attached to the phone's audio port.
	AssignLight(ctx context.Context, request *AssignLightRequest) (response *AssignLightResponse, err error)

	// ControllerName returns the name of the controller.
	ControllerName(ctx context.Context, request *ControllerNameRequest) (response *ControllerNameResponse, err error)

	// ExtinguishAll sets all group's to 0% intensity, toggles all themes
	// to "off", and disables FlashLights mode.
	ExtinguishAll(ctx context.Context, request *ExtinguishAllRequest) (response *ExtinguishAllResponse, err error)

	// FlashLights enters or leaves a special mode used while assigning
	// lights to groups. (Possibly "flash" refers to writing to flash
	// rather than a light shining brightly and briefly. The protocol for
	// learning serial numbers and/or the AssignLight method may only work
	// in this mode.) In this mode, all lights are on at 100%, which is
	// not reflected in GroupListGet's returned intensities. On leaving,
	// lights turn off and all intensities are set to 0%.
	FlashLights(ctx context.Context, request *FlashLightsRequest) (response *FlashLightsResponse, err error)

	// GroupListAdd adds a group, appending it to the list.
	GroupListAdd(ctx context.Context, request *GroupListAddRequest) (response *GroupListAddResponse, err error)

	// GroupListClear deletes all groups.
	GroupListClear(ctx context.Context, request *GroupListClearRequest) (response *GroupListClearResponse, err error)

	// GroupListDelete deletes a group.
	GroupListDelete(ctx context.Context, request *GroupListDeleteRequest) (response *GroupListDeleteResponse, err error)

	// GroupListGet retrieves all groups, including their current intensity.
	GroupListGet(ctx context.Context, request *GroupListGetRequest) (response *GroupListGetResponse, err error)

	// GroupListRename renames a group.
	GroupListRename(ctx context.Context, request *GroupListRenameRequest) (response *GroupListRenameResponse, err error)

	// GroupListReorder reorders groups.
	GroupListReorder(ctx context.Context, request *GroupListReorderRequest) (response *GroupListReorderResponse, err error)

	// IlluminateAll illuminates all lights to 75% intensity.
	IlluminateAll(ctx context.Context, request *IlluminateAllRequest) (response *IlluminateAllResponse, err error)

	// IlluminateGroup illuminates a single group to given intensity,
	// without affecting other groups.
	IlluminateGroup(ctx context.Context, request *IlluminateGroupRequest) (response *IlluminateGroupResponse, err error)

	// IlluminateTheme illuminates all groups in a theme to either the
	// intensity defined for them in the theme, or 0%.
	IlluminateTheme(ctx context.Context, request *IlluminateThemeRequest) (response *IlluminateThemeResponse, err error)

	// ThemeClear removes all groups from a theme.
	ThemeClear(ctx context.Context, request *ThemeClearRequest) (response *ThemeClearResponse, err error)

	// ThemeGet retrieves a theme's definition, including its name and (group, intensity) list.
	ThemeGet(ctx context.Context, request *ThemeGetRequest) (response *ThemeGetResponse, err error)

	// ThemeListAdd adds a new theme.
	ThemeListAdd(ctx context.Context, request *ThemeListAddRequest) (response *ThemeListAddResponse, err error)

	// ThemeListClear deletes all themes.
	ThemeListClear(ctx context.Context, request *ThemeListClearRequest) (response *ThemeListClearResponse, err error)

	// ThemeListDelete deletes one theme.
	ThemeListDelete(ctx context.Context, request *ThemeListDeleteRequest) (response *ThemeListDeleteResponse, err error)

	// ThemeListGet retrieves the list of themes and their current status.
	// Does not include the full definition of each theme.
	ThemeListGet(ctx context.Context, request *ThemeListGetRequest) (response *ThemeListGetResponse, err error)

	// ThemeListRename renames a theme.
	ThemeListRename(ctx context.Context, request *ThemeListRenameRequest) (response *ThemeListRenameResponse, err error)

	// ThemeListReorder reorders the list of themes.
	ThemeListReorder(ctx context.Context, request *ThemeListReorderRequest) (response *ThemeListReorderResponse, err error)

	// ThemeSet redefines a theme.
	ThemeSet(ctx context.Context, request *ThemeSetRequest) (response *ThemeSetResponse, err error)
}

const (
	MaxIntensity   = 100
	MaxThemeNumber = 25
	MaxNameLength  = 19

	StatusOk                   = 0
	StatusUnknownMethod        = 1
	StatusUnparseableRequest   = 101
	StatusInvalidRequest       = 102
	StatusPreconditionFailed   = 201
	StatusGroupNameInUse       = 202
	StatusGroupNumberInUse     = 205
	StatusThemeIndexOutOfRange = 243
)

var statusName = map[int]string{
	StatusOk:                   "ok",
	StatusUnknownMethod:        "unknown method",
	StatusUnparseableRequest:   "unparseable request",
	StatusInvalidRequest:       "invalid request",
	StatusPreconditionFailed:   "precondition failed",
	StatusGroupNameInUse:       "group name in use",
	StatusGroupNumberInUse:     "group number in use",
	StatusThemeIndexOutOfRange: "theme index out of range",
}

// ErrorForStatus returns an error for the given status.
// If status == 0, the error will be nil.
func ErrorForStatus(status int) error {
	if status == 0 {
		return nil
	}
	if statusStr, ok := statusName[status]; ok {
		return errors.New(statusStr)
	}
	return fmt.Errorf("unknown status %d", status)
}

type AssignLightRequest struct {
	SerialNumber int
	GroupNumber  uint8
}

type AssignLightResponse struct {
	Status int
}

type ControllerNameRequest struct {
}

type ControllerNameResponse struct {
	Status     int
	Controller string
}

type ExtinguishAllRequest struct {
}

type ExtinguishAllResponse struct {
	Status int
}

type FlashLightsRequest struct {
	OnOff uint8
}

type FlashLightsResponse struct {
	Status int
}

type Group struct {
	GroupNumber uint8
	Intensity   uint8
	Name        string
}

type GroupListAddRequest struct {
	GroupNumber uint8

	// Name will be truncated to MaxNameLength.
	Name string
}

type GroupListAddResponse struct {
	// Status will be StatusGroupNumberInUse if the number is already
	// taken, or StatusGroupNameInUse if the name is already taken.
	Status int
}

type GroupListClearRequest struct {
}

type GroupListClearResponse struct {
	Status    int
	GroupList []Group
}

type GroupListDeleteRequest struct {
	Name string
}

type GroupListDeleteResponse struct {
	// Status will be StatusPreconditionFailed if the group does not
	// exist.
	Status int
}

type GroupListGetRequest struct {
}

type GroupListGetResponse struct {
	Status    int
	GroupList []Group
}

type GroupListRenameRequest struct {
	OldName string

	// NewName will be truncated to MaxNameLength.
	NewName string
}

type GroupListRenameResponse struct {
	// Status will be StatusGroupNameInUse if the name is taken.
	Status int
}

type GroupListReorderRequest struct {
	// GroupNumbers should be a new order that includes all existing
	// groups exactly once.
	GroupNumbers []uint8
}

type GroupListReorderResponse struct {
	// Status will be StatusPreconditionFailed if the request did not list
	// all group numbers exactly once.
	Status int
}

type IlluminateAllRequest struct {
}

type IlluminateAllResponse struct {
	Status int
}

type IlluminateGroupRequest struct {
	GroupNumber uint8
	Intensity   uint8
}

type IlluminateGroupResponse struct {
	Status int
}

type IlluminateThemeRequest struct {
	ThemeIndex uint8

	// OnOff should be 0 to set the intensity to 0, or non-zero to
	// use the intensities stored in the theme.
	OnOff uint8
}

type IlluminateThemeResponse struct {
	Status int
}

type Theme struct {
	Name       string
	ThemeIndex uint8
	OnOff      uint8
}

type ThemeGroup struct {
	GroupNumber uint8
	Intensity   uint8
}

type ThemeClearRequest struct {
	ThemeIndex uint8
}

type ThemeClearResponse struct {
	Status int
}

type ThemeGetRequest struct {
	ThemeIndex uint8
}

type ThemeGetResponse struct {
	Status int
	Groups []ThemeGroup
}

type ThemeListAddRequest struct {
	ThemeIndex uint8

	// Name will be truncated to MaxNameLength.
	Name string
}

type ThemeListAddResponse struct {
	Status int
}

type ThemeListClearRequest struct {
}

type ThemeListClearResponse struct {
	Status int
}

type ThemeListDeleteRequest struct {
	Name string
}

type ThemeListDeleteResponse struct {
	// Status will be StatusInvalidRequest if the themes are restricted,
	// or StatusPreconditionFailed if the theme does not exist.
	Status int
}

type ThemeListGetRequest struct {
}

type ThemeListGetResponse struct {
	Status int

	// Restricted will be non-zero iff themes are restricted in the
	// controller's setup menu. In this case, some operations on themes
	// can not be performed through this API.
	Restricted int

	ThemeList []Theme
}

type ThemeListRenameRequest struct {
	OldName string

	// NewName will be truncated to MaxNameLength.
	NewName string
}

type ThemeListRenameResponse struct {
	Status int
}

type ThemeListReorderRequest struct {
	ThemeIndexes []uint8
}

type ThemeListReorderResponse struct {
	// Status will be StatusPreconditionFailed if the request did not list
	// all group numbers exactly once.
	Status int
}

type ThemeSetRequest struct {
	ThemeIndex uint8
	Groups     []ThemeGroup
}

type ThemeSetResponse struct {
	Status int
}
