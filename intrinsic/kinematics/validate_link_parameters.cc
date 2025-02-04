// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/kinematics/validate_link_parameters.h"

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/math/inertia_utils.h"

namespace intrinsic::kinematics {

icon::RealtimeStatus ValidateMass(double mass_kg) {
  if (mass_kg <= 0) {
    return icon::FailedPreconditionError(icon::RealtimeStatus::StrCat(
        "The mass should be > 0.0, but got ", mass_kg, " kg instead."));
  }
  return icon::OkStatus();
}

icon::RealtimeStatus ValidateInertia(const eigenmath::Matrix3d& inertia) {
  return intrinsic::ValidateInertia(inertia);
}

}  // namespace intrinsic::kinematics
