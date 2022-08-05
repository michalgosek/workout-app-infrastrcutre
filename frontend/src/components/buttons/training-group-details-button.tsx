import { FC, PropsWithChildren } from 'react';

import { Button } from "react-bootstrap";
import style from './training-group-details.module.scss';

type TrainingGroupDetailsButtonProps = {
    trainingUUID: string;
    trainerUUID: string;
};

const TrainerGroupDetailsButton: FC<PropsWithChildren<TrainingGroupDetailsButtonProps>> = ({ trainingUUID, trainerUUID }) => {
    const TRAINING_GROUP_DETAILS_ENDPOINT = `http://localhost:3000/trainers/${trainerUUID}/trainings/${trainingUUID}`;
    return <Button type="button" className={style.detailsbutton} href={TRAINING_GROUP_DETAILS_ENDPOINT}>Details</Button>
}

export default TrainerGroupDetailsButton;