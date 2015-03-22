package service

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#import <IOKit/IOKitLib.h>
#import <CoreFoundation/CoreFoundation.h>
#include <errno.h>

static io_connect_t dataPort = 0;
struct reading {
   uint64_t left;
   uint64_t right;
};

struct reading getReading() {
  struct reading out;
  io_service_t serviceObject;
  serviceObject = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching("AppleLMUController"));
  if (!serviceObject) {
    errno = 100;
    return out;
  }

  errno = IOServiceOpen(serviceObject, mach_task_self(), 0, &dataPort);
  IOObjectRelease(serviceObject);
  if (errno != KERN_SUCCESS) {
    return out;
  }
  uint32_t outputs = 2;
  uint64_t values[outputs];

  errno = IOConnectCallMethod(dataPort, 0, nil, 0, nil, 0, values, &outputs, nil, 0);
  if (errno == KERN_SUCCESS) {
    out.left = values[0];
    out.right = values[1];
  }
  return out;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"syscall"
)

var ErrNoLightSensors = errors.New("your machine doesn't have an ambient light sensor")

type Reading struct {
	Left  int64
	Right int64
	Mean  int64
}

func (r Reading) String() string {
	return fmt.Sprintf("left: %d right: %d mean: %d", r.Left, r.Right, r.Mean)
}

func readSensors() (Reading, error) {
	var (
		reading Reading
		err     error
	)
	rc, errno := C.getReading()
	if err != nil {
		switch errno.(syscall.Errno) {
		case 100:
			err = errors.New("your machine doesn't have an ambient light sensor")
		default:
			err = errors.New("problem getting reading")
		}
	}

	reading.Left, reading.Right = int64(rc.left), int64(rc.right)
	reading.Mean = (reading.Left + reading.Right) / 2
	return reading, err
}
