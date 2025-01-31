# Copyright 2023 Intrinsic Innovation LLC

"""A skill that echos its parameters as its result."""

from intrinsic.skills.python import skill_interface as skl
from intrinsic.skills.python import skill_interface_utils
from intrinsic.skills.testing import echo_skill_pb2
from intrinsic.util import decorators


class EchoSkill(skl.Skill):
  """A skill that echos its parameters as its result."""

  @decorators.overrides(skl.Skill)
  def execute(
      self,
      request: skl.ExecuteRequest[echo_skill_pb2.EchoSkillParams],
      context: skl.ExecuteContext,
  ) -> echo_skill_pb2.EchoSkillReturn:
    return echo_skill_pb2.EchoSkillReturn(foo=request.params.foo)

  @decorators.overrides(skl.Skill)
  def preview(
      self,
      request: skl.PreviewRequest[echo_skill_pb2.EchoSkillParams],
      context: skl.PreviewContext,
  ) -> echo_skill_pb2.EchoSkillReturn:
    return skill_interface_utils.preview_via_execute(self, request, context)
