from typing import Set

import pytest

from abi.binary.executable import Executable
from abi.section.section_flag import SectionFlag
from abi.section.section_name import SectionName
from abi.section.section_type import SectionType
from initializer.int_initializer import IntInitializer
from initializer.str_initializer import StrInitializer


@pytest.fixture
def executable() -> Executable:
    return Executable("", set())


@pytest.fixture
def section_name() -> SectionName:
    return SectionName.TEXT


@pytest.fixture
def name() -> str:
    return StrInitializer.identifier(IntInitializer.value_by_range(1, 64))


@pytest.fixture
def section_flags() -> Set[SectionFlag]:
    return {SectionFlag.ALLOC, SectionFlag.EXECINSTR}


@pytest.fixture
def section_type() -> SectionType:
    return SectionType.PROG_BITS


def test_checkout_section(
    executable: Executable,
    section_name: SectionName,
    name: str,
    section_flags: Set[SectionFlag],
    section_type: SectionType,
):
    executable.checkout_section(section_name, name, section_flags, section_type)
    assert executable.section(section_name, name) is not None
