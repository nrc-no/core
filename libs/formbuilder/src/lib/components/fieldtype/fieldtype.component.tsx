import React, { Fragment, FC } from 'react';
import {
    FormLabel,
    FormCheck,
    FormCheckLabel,
    FormRadioInput
} from '@nrc.no/ui-toolkit';

type FieldTypeProps = {
    value: FieldType | undefined
}

export enum FieldType {
    date = "Date",
    text = "Text",
    textarea = "Textarea",
    integer = "Integer",
    checkbox = "Checkbox",
    dropdown = "Dropdown",
    multidropdown = "Dropdown (multi)"
}

const makeFieldTypeRadios = (selected: FieldType | undefined) => {
    const returnList = []
    for (const option in FieldType){
        if (selected === FieldType[option]) {
            returnList.push(
                <FormCheck>
                    <FormCheckLabel>{FieldType[option]}</FormCheckLabel>
                    <FormRadioInput name={'fieldtype'} value={FieldType[option]} checked={true}/>
                </FormCheck>
            )
        } else {
            returnList.push(
                <FormCheck>
                    <FormCheckLabel>{FieldType[option]}</FormCheckLabel>
                    <FormRadioInput name={'fieldtype'} value={FieldType[option]}/>
                </FormCheck>
            )
        }
    }
    return returnList
}

export const FieldTypePicker: FC<FieldTypeProps> = (props) => {
    return (
        <Fragment>
            {
                makeFieldTypeRadios(props.value).map((radio) => radio)
            }
        </Fragment>
    )
}
