// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/world/robot_payload/robot_payload_base.h"

#include <ostream>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/almost_equals.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/world/proto/robot_payload.pb.h"

namespace intrinsic {

RobotPayloadBase::RobotPayloadBase()
    : mass_kg_(0.0),
      tip_t_cog_(Pose3d::Identity()),
      inertia_in_cog_(eigenmath::Matrix3d::Zero()) {}

RobotPayloadBase::RobotPayloadBase(double mass, const Pose3d& tip_t_cog,
                                   const eigenmath::Matrix3d& inertia)
    : mass_kg_(mass), tip_t_cog_(tip_t_cog), inertia_in_cog_(inertia) {}

bool RobotPayloadBase::operator==(const RobotPayloadBase& other) const {
  return IsApprox(other, kStdError);
}

std::ostream& operator<<(std::ostream& os, const RobotPayloadBase& payload) {
  os << "Payload: mass: " << payload.mass()
     << " tip_t_cog: " << payload.tip_t_cog()
     << " inertia: " << payload.inertia();
  return os;
}

bool RobotPayloadBase::IsApprox(const intrinsic::RobotPayloadBase& other,
                                double precision) const {
  return intrinsic::AlmostEquals(mass(), other.mass(), precision) &&
         tip_t_cog().isApprox(other.tip_t_cog(), precision) &&
         inertia().isApprox(other.inertia(), precision);
}
}  // namespace intrinsic
