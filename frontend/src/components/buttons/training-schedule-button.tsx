import { FC, PropsWithChildren } from 'react';

import { Button } from 'react-bootstrap';
import { TrainingsService } from 'services/trainings-service';
import { useAuth0 } from '@auth0/auth0-react';

type TrainingScheduleProps = {
    trainingUUID: string;
    trainerUUID: string;
};

const isValidTrainingScheduleParams = (trainerUUID: string, trainingUUID: string, userUUID: string, userName: string): boolean => {
    const params: string[] = [
        trainerUUID,
        trainingUUID,
        userUUID,
        userName
    ]
    return !params.some(p => p === "");
}

const TrainingScheduleButton: FC<PropsWithChildren<TrainingScheduleProps>> = ({ trainerUUID, trainingUUID }) => {
    const { isAuthenticated, user, loginWithRedirect, } = useAuth0();
    const AUTH0_CALLBACK_URL: string = process.env.REACT_APP_AUTH0_CALLBACK_URL || "";

    const onClickHandle = () => {
        if (!isAuthenticated) {
            loginWithRedirect({
                redirectUri: AUTH0_CALLBACK_URL,
            }).then().catch(err => console.error(err));
        }

        const userUUID = user?.sub ?? "";
        const userName = user?.name ?? "";
        const isParamsValid = isValidTrainingScheduleParams(trainerUUID, trainingUUID, userUUID, userName)
        if (!isParamsValid) return;

        const isTrainerTrainingGroup = trainerUUID === user?.sub
        if (isTrainerTrainingGroup) return;

        const signUpParticipant = async () => {
            const res = await TrainingsService.signupParticipant({
                trainer_uuid: trainerUUID,
                trainer_group_uuid: trainingUUID,
                participant: {
                    participant_uuid: userUUID,
                    participant_name: userName,
                }
            });
            console.log(res);
        }
        signUpParticipant();
    };
    return (
        <Button type="button" className="btn btn-primary" onClick={onClickHandle}>Schedule</Button>
    );
};

export default TrainingScheduleButton;


