from typing import Set

import pytest

from abi.section.section import Section
from abi.section.section_flag import SectionFlag
from abi.section.section_name import SectionName
from abi.section.section_type import SectionType
from abi.word.data_address_word import DataAddressWord
from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from initializer.int_initializer import IntInitializer
from initializer.str_initializer import StrInitializer


@pytest.fixture
def section_names() -> Set[SectionName]:
    return set(SectionName)


@pytest.fixture
def name() -> str:
    return StrInitializer.identifier(IntInitializer.value_by_range(1, 64))


@pytest.fixture
def section_flags() -> Set[SectionFlag]:
    return {SectionFlag.ALLOC, SectionFlag.EXECINSTR}


@pytest.fixture
def section_type() -> SectionType:
    return SectionType.PROG_BITS


def test_address(
    section_names: Set[SectionName],
    section_flags: Set[SectionFlag],
    section_type: SectionType,
):
    for section_name in section_names:
        for _ in range(100):
            address = (
                IntInitializer.value_by_width(Representation.UNSIGNED, DataAddressWord().width())
                // InstructionWord().size()
            ) * InstructionWord().size()

            section = Section(section_name, "", section_flags, section_type)
            section.set_address(address)

            assert address == section.address()


def test_label(
    section_names: Set[SectionName],
    name: str,
    section_flags: Set[SectionFlag],
    section_type: SectionType,
):
    for section_name in section_names:
        for _ in range(100):
            section = Section(section_name, name, section_flags, section_type)
            section.append_label(name)

            label = section.label(name)
            assert label is not None
            assert label.name() == name
