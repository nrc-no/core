import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M13,17v-4h4v4H13z M11,11.5c0-0.3,0.2-0.5,0.5-0.5h7c0.3,0,0.5,0.2,0.5,0.5v7c0,0.3-0.2,0.5-0.5,0.5h-7
	c-0.3,0-0.5-0.2-0.5-0.5V11.5z M13,27v-4h4v4H13z M11,21.5c0-0.3,0.2-0.5,0.5-0.5h7c0.3,0,0.5,0.2,0.5,0.5v7c0,0.3-0.2,0.5-0.5,0.5
	h-7c-0.3,0-0.5-0.2-0.5-0.5V21.5z M23,13v4h4v-4H23z M21.5,11c-0.3,0-0.5,0.2-0.5,0.5v7c0,0.3,0.2,0.5,0.5,0.5h7
	c0.3,0,0.5-0.2,0.5-0.5v-7c0-0.3-0.2-0.5-0.5-0.5H21.5z M23,27v-4h4v4H23z M21,21.5c0-0.3,0.2-0.5,0.5-0.5h7c0.3,0,0.5,0.2,0.5,0.5
	v7c0,0.3-0.2,0.5-0.5,0.5h-7c-0.3,0-0.5-0.2-0.5-0.5V21.5z"
        />
    )
}
