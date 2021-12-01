import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M15.1,18.7c-0.2-0.2-0.3-0.5-0.3-0.8c0-0.3,0.1-0.6,0.3-0.8c0.2-0.2,0.5-0.3,0.8-0.3c0.3,0,0.6,0.1,0.8,0.3
	c0.2,0.2,0.3,0.5,0.3,0.8c0,0.3-0.1,0.6-0.3,0.8c-0.2,0.2-0.5,0.3-0.8,0.3C15.6,19,15.3,18.9,15.1,18.7z M23.3,18.7
	c-0.2-0.2-0.3-0.5-0.3-0.8c0-0.3,0.1-0.6,0.3-0.8c0.2-0.2,0.5-0.3,0.8-0.3s0.6,0.1,0.8,0.3c0.2,0.2,0.3,0.5,0.3,0.8
	c0,0.3-0.1,0.6-0.3,0.8c-0.2,0.2-0.5,0.3-0.8,0.3S23.5,18.9,23.3,18.7z M9.5,20c0-5.8,4.7-10.5,10.5-10.5S30.5,14.2,30.5,20
	S25.8,30.5,20,30.5S9.5,25.8,9.5,20z M23.4,28c1-0.4,2-1.1,2.8-1.9c0.8-0.8,1.4-1.7,1.9-2.8c0.5-1.1,0.7-2.2,0.7-3.4
	s-0.2-2.3-0.7-3.4c-0.4-1-1.1-2-1.9-2.8c-0.8-0.8-1.7-1.4-2.8-1.9c-1.1-0.5-2.2-0.7-3.4-0.7s-2.3,0.2-3.4,0.7c-1,0.4-2,1.1-2.8,1.9
	c-0.8,0.8-1.4,1.7-1.9,2.8c-0.5,1.1-0.7,2.2-0.7,3.4s0.2,2.3,0.7,3.4c0.4,1,1.1,2,1.9,2.8c0.8,0.8,1.7,1.4,2.8,1.9
	c1.1,0.5,2.2,0.7,3.4,0.7S22.3,28.5,23.4,28z M21.4,23c0,0.8-0.2,1.4-0.5,1.8c-0.3,0.4-0.6,0.5-0.8,0.5c-0.2,0-0.5-0.1-0.8-0.5
	c-0.3-0.4-0.5-1.1-0.5-1.8c0-0.8,0.2-1.4,0.5-1.8c0.3-0.4,0.6-0.5,0.8-0.5c0.2,0,0.5,0.1,0.8,0.5C21.1,21.6,21.4,22.2,21.4,23z
	 M22.7,23c0,2-1.2,3.7-2.7,3.7c-1.5,0-2.7-1.6-2.7-3.7c0-2,1.2-3.7,2.7-3.7C21.5,19.3,22.7,21,22.7,23z"
        />
    )
}
