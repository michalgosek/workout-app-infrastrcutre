import { FC, PropsWithChildren } from 'react';

import { Button } from 'react-bootstrap';
import { TrainingsService } from 'services/trainings-service';
import { useAuth0 } from '@auth0/auth0-react';

type TrainingScheduleProps = {
    trainingUUID: string;
    trainerUUID: string;
};

const TrainingScheduleButton: FC<PropsWithChildren<TrainingScheduleProps>> = ({ trainerUUID, trainingUUID }) => {
    const { isAuthenticated, loginWithRedirect, user } = useAuth0();
    const AUTH0_CALLBACK_URL: string = process.env.REACT_APP_AUTH0_CALLBACK_URL || "";

    const onClickHandle = () => {
        if (!isAuthenticated) {
            loginWithRedirect({ redirectUri: AUTH0_CALLBACK_URL })
                .then()
                .catch(err => console.error(err));
        }
        if (!user) return null;

        TrainingsService.signupParticipant({
            trainer_uuid: trainerUUID,
            trainer_group_uuid: trainingUUID,
            participant: {
                participant_uuid: user.sub ?? '',
                participant_name: user.name ?? ''
            }
        })
            .then(res => console.log(res))
            .catch(err => console.error(err));

        window.location.reload();
    };
    return (
        <Button type="button" className="btn btn-primary" onClick={onClickHandle}>Schedule</Button>
    );
};

export default TrainingScheduleButton;


