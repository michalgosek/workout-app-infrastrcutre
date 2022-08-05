import React, { FC } from 'react';

import TrainingForm from './components/form';
import { TrainingsService } from 'services/trainings-service';

const TrainingPlanningForm: FC = () => {
    return (
        <TrainingForm placeholders={
            {
                description: "this is example description",
                name: "this is example name",
                appointment: new Date().toString(),
            }
        } callbackAPI={TrainingsService.createTrainingGroup} />
    )
};

const PlanTraining: React.FC = () => {
    return (
        <>
            <TrainingPlanningForm />
        </>
    );
};

export default PlanTraining;