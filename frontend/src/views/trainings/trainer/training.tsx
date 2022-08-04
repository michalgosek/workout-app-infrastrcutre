import { FC, PropsWithChildren } from 'react';

type TrainerTrainingGroupProps = {
    trainerUUID: string;
    trainingUUID: string;
};

const TrainerTrainingGroup: FC<PropsWithChildren<TrainerTrainingGroupProps>> = ({ trainerUUID, trainingUUID }) => {
    const details = `TrainigUUID:${trainingUUID}, TrainerUUID: ${trainerUUID}`;
    return (
        <h1>Training Details for:`${details}`</h1>
    );
};


export default TrainerTrainingGroup;