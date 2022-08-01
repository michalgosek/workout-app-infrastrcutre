import { PlanTrainingGroupPOST, TrainingsService } from "../../../services/trainings-service";
import React, { FC, Fragment, useState } from "react";
import { SubmitHandler, useForm } from 'react-hook-form';

import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import Button from 'react-bootstrap/Button';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import Form from 'react-bootstrap/Form';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import { useAuth0 } from "@auth0/auth0-react";

type FormValues = {
    groupName: string;
    appointment: Date;
    description: string;
};

const TrainingPlanningForm: FC = () => {
    const [value, setValue] = useState<Date | null>(new Date());
    const { handleSubmit, register, formState: { errors } } = useForm<FormValues>();
    const { user } = useAuth0();
    if (!user) {
        return null;
    }

    const onHandleChange = (newValue: Date | null) => {
        setValue(newValue);
    };
    const onSubmitHandle: SubmitHandler<FormValues> = (data) => {
        (async () => {
            const name = user.name ?? "not_defined";
            const uuid = user.sub ?? "not_defined";
            const training: PlanTrainingGroupPOST = {
                user: {
                    role: "Trainer",
                    uuid: uuid,
                    name: name,
                },
                group_name: data.groupName,
                group_desc: data.description,
                date: data.appointment.toString(),
            }
            await TrainingsService.createTrainingGroup(training);
            window.location.reload();
        })();
    }
    return (
        <Form onSubmit={handleSubmit(onSubmitHandle)} className="form">
            <Form.Group className="mb-3" controlId="formBasicEmail">
                <Form.Label>Training group name:</Form.Label>
                <Form.Control
                    type="text"
                    placeholder="Training group name"
                    {...register("groupName", { required: "training group name is obligatory!" })}
                />
                {errors?.groupName && <p>{errors.groupName.message}</p>}
                <Form.Text className="text-muted">
                    This will be used as group name description.
                </Form.Text>
            </Form.Group>
            <Form.Group className="mb-3" controlId="formBasicDate">
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                    <Stack spacing={3}>
                        <DateTimePicker
                            {...register('appointment', { required: "appointment date for training group is obligatory!" })}
                            label="Training appointment"
                            value={value}
                            onChange={onHandleChange}
                            ampm={false}
                            renderInput={(params) => <TextField {...params} />}
                        />
                    </Stack>
                </LocalizationProvider>
            </Form.Group>
            <Form.Group className="mb-3" controlId="formBasicTextArea">
                <Form.Label>Training group description:</Form.Label>
                <Form.Control
                    as="textarea"
                    placeholder="Leave a example description here"
                    {...register("description", { required: "This value is obligatory!" })}
                    style={{ height: '100px' }}
                />
            </Form.Group>
            <Button className="btn btn-primary me-2" type="submit">Plan</Button>
            <Button className="btn btn-secondary" type="reset" >Clear</Button>
        </Form >
    );
};



const PlanTraining: React.FC = () => {
    return (
        <Fragment>
            <TrainingPlanningForm />
        </Fragment>
    );
};

export default PlanTraining;