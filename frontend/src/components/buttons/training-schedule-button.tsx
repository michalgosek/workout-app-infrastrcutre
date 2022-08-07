import { FC, PropsWithChildren } from 'react';

import { Button } from 'react-bootstrap';
import { TrainingsService } from 'services/trainings-service';
import { useAuth0 } from '@auth0/auth0-react';
import { useGetTrainingsServiceAccessToken } from 'services/hooks';

type TrainingScheduleProps = {
    trainingUUID: string;
    trainerUUID: string;
};

const TrainingScheduleButton: FC<PropsWithChildren<TrainingScheduleProps>> = ({ trainerUUID, trainingUUID }) => {
    const { isAuthenticated, loginWithRedirect, user } = useAuth0();
    const token = useGetTrainingsServiceAccessToken();
    const AUTH0_CALLBACK_URL: string = process.env.REACT_APP_AUTH0_CALLBACK_URL || "";

    const onClickHandle = () => {
        const callRedirectAPI = async (isAuthenticated: boolean) => {
            if (isAuthenticated) {
                console.log('user already authenitcated');
                return;
            }
            try {
                const res = await loginWithRedirect({ redirectUri: AUTH0_CALLBACK_URL })
                console.log(res);
            } catch (err) {
                console.error(err);
            }
        }

        const signupParticipant = async () => {
            if (!user) return null;

            const res = await TrainingsService.signupParticipant({
                trainer_uuid: trainerUUID,
                trainer_group_uuid: trainingUUID,
                participant: {
                    participant_uuid: user.sub ?? '',
                    participant_name: user.name ?? ''
                }
            }, token)
            console.log(res);
            window.location.reload();
        }
        callRedirectAPI(isAuthenticated);
        signupParticipant();
    };
    return (
        <Button type="button" className="btn btn-primary" onClick={onClickHandle}>Schedule</Button>
    );
};

export default TrainingScheduleButton;


