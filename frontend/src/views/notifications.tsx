import { Alert, Container } from 'react-bootstrap';
import notificationService, { NotificationReadModel } from '../services/notification-service';
import { useEffect, useState } from "react";

import DropButton from 'components/buttons/drop-button';
import style from './notifications.module.scss';
import { useAuth0 } from '@auth0/auth0-react';
import { useParams } from 'react-router-dom';

const dropAllNotifications = (userUUID: string, token: string) => {
    const deleteAll = async () => {
        const res = await notificationService.deleteNotifications(userUUID, token)
        console.log(res);
        window.location.reload();
    }
    deleteAll();
}

const Notifications = () => {
    const [notifications, setNotifications] = useState<NotificationReadModel[]>();
    const { getAccessTokenSilently } = useAuth0();
    const [token, setToken] = useState('');

    const { participantUUID } = useParams();

    useEffect(() => {
        if (!participantUUID) {
            console.error('missing participantUUID in URL path');
            return;
        }
        getAccessTokenSilently().then(token => {
            setToken(token)
            notificationService
                .getAllNotifications(participantUUID, token)
                .then((notifications) => setNotifications(notifications))
                .catch(err => console.error(err))
        }).catch(err => console.error(err));
    }, [getAccessTokenSilently, participantUUID]);


    if (!participantUUID) {
        console.error('missing participantUUID in URL path');
        return null;
    }


    return (
        <>
            <Container>
                {
                    notifications?.map(n => (
                        <Alert key={n.training_uuid} variant='info'>
                            <Alert.Heading>{n.title}</Alert.Heading>
                            <hr />
                            <div>
                                <p><b>Date:</b> {n.date}</p>
                                <p><b>Message:</b> {n.content}</p>
                                <p><b>Trainer:</b> {n.trainer}</p>
                            </div>
                        </Alert>
                    ))
                }
            </Container>
            {
                (notifications?.length !== 0) ?
                    <Container className={style["vertical-center"]}>
                        <DropButton onClickHandle={() => dropAllNotifications(participantUUID, token)} />
                    </Container>
                    : <h1>There is no notifications!  </h1>
            }
        </>

    );
};

export default Notifications; 