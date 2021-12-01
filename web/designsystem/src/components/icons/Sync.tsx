import React from "react";
import {Path} from "react-native-svg"
import {IconVariants} from "../../types/icons";

export default (variant: IconVariants) => {
    return (
        <Path
            fill={variant}
            d="M11.9,19.8c0-1,0.2-2,0.6-3c0.4-1,1-1.8,1.7-2.6s1.6-1.3,2.6-1.7c1-0.4,2-0.6,3.1-0.6c1.1,0,2.1,0.2,3.1,0.6
	c1,0.4,1.8,1,2.6,1.7c0.2,0.2,0.4,0.5,0.7,0.7l-1.4,1.1c0,0,0,0.1-0.1,0.1c0,0,0,0.1,0,0.1c0,0,0,0.1,0,0.1c0,0,0.1,0,0.1,0.1l4.1,1
	c0.1,0,0.2-0.1,0.2-0.2l0-4.2c0-0.2-0.2-0.2-0.3-0.1l-1.3,1c-1.8-2.3-4.6-3.8-7.8-3.8c-5.4,0-9.7,4.3-9.8,9.7c0,0,0,0,0,0.1
	c0,0,0,0,0,0.1c0,0,0,0,0.1,0c0,0,0,0,0.1,0h1.4C11.9,20,11.9,19.9,11.9,19.8z M29.7,20h-1.4c-0.1,0-0.2,0.1-0.2,0.2
	c0,1-0.2,2-0.6,3c-0.4,1-1,1.8-1.7,2.6c-0.7,0.8-1.6,1.3-2.6,1.8s-2,0.6-3.1,0.6c-1.1,0-2.1-0.2-3.1-0.6c-1-0.4-1.9-1-2.6-1.8
	c-0.2-0.2-0.4-0.5-0.7-0.7l1.4-1.1c0,0,0-0.1,0.1-0.1c0,0,0-0.1,0-0.1c0,0,0-0.1,0-0.1c0,0-0.1,0-0.1-0.1l-4.1-1
	c-0.1,0-0.2,0.1-0.2,0.2l0,4.2c0,0.2,0.2,0.2,0.3,0.1l1.3-1c1.8,2.3,4.6,3.8,7.8,3.8c5.4,0,9.7-4.3,9.8-9.7c0,0,0,0,0-0.1
	c0,0,0,0,0-0.1C29.8,20,29.8,20,29.7,20C29.7,20,29.7,20,29.7,20z"
        />
    )
}
