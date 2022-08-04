import {FC, MouseEventHandler, PropsWithChildren} from 'react';

import { Button } from "react-bootstrap";

interface DropButtonProps {
    onClickHandle: MouseEventHandler;
}

const DropButton: FC<PropsWithChildren<DropButtonProps>> = ({onClickHandle}) => {
    return  (
        <Button type="button" className="btn btn-danger" onClick={onClickHandle}>Drop</Button>
    );
};

export default DropButton;


 