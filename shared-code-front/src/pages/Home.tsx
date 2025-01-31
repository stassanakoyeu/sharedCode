import { TextField, Typography, useTheme, Grid, InputAdornment } from "@mui/material"
import React, { useState } from "react";
import CreateRoomNamePrompt from "../components/CreateRoom_NamePrompt";
import JoinRoomNamePrompt from "../components/JoinRoom_NamePrompt";
import CodeSvgDark from "../assets/undraw_collaboration_dark.png";
import CodeSvgLight from "../assets/undraw_collaboration_light.png";
import PrimaryButton from "../components/ui/PrimaryButton";
import SecondaryButton from "../components/ui/SecondaryButton";
import MeetingRoomIcon from '@mui/icons-material/MeetingRoom';
import KeyboardIcon from '@mui/icons-material/Keyboard';
import { getRoomById } from "../services/room";


export default function Home() {
    const theme = useTheme();
    
    const [roomId, setRoomId] = useState("");
    const [createRoomOpen, setCreateRoomOpen] = useState(false);
    const [joinRoomOpen, setJoinRoomOpen] = useState(false);

    const [joinError, setJoinError] = useState("");

    const handleRoomIdChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setRoomId(e.target.value);
        setJoinError("");
    }

    const handleJoinRoom = () => {
        getRoomById(roomId)
        .then((data) => {
            if(data.status === 404) {
                setJoinError(data.error);
            }
            else {
                setJoinError("");
                setJoinRoomOpen(true);
            }
        })
    }

    const collabSvg = () => {
        if(theme.palette.mode == 'dark') {
            return (<img style={{ height: "60%", width: "60%" }} src={CodeSvgDark} alt="Code Companion Illustration" />              );
        }
        else {
            return (<img style={{ height: "60%", width: "60%" }} src={CodeSvgLight} alt="Code Companion Illustration" />              );
        }
    }
    return (
        <>
            <Grid container maxWidth="lg" sx={{ height: "100%", alignItems: 'center' }}>
                <Grid item xs={12} md={6} sx={{ paddingInline: '1em', paddingBlock: '1em' }}>
                    <Typography variant="h2" marginBottom={3}>Платформа для совместного написания кода.<br /></Typography>
                    <Typography variant="body1" marginBottom={3}>Распределенный онлайн-редактор с возможность подсветки кода. Лучшее решение для интервью, коллаборативного программирования и учёбы</Typography>
                        <PrimaryButton size="large" onClick={() => setCreateRoomOpen(true)} startIcon={<MeetingRoomIcon />}>
                            Создать комнату
                        </PrimaryButton>
                        <CreateRoomNamePrompt 
                         open={createRoomOpen} 
                         setOpen={setCreateRoomOpen} 
                        />
                        <JoinRoomNamePrompt
                         open={joinRoomOpen} 
                         setOpen={setJoinRoomOpen}
                         roomId={roomId}
                        />
                        &nbsp;
                        или
                        &nbsp;
                        <TextField
                            error={joinError !== ""}
                            helperText={joinError}
                            InputProps={{
                                startAdornment: (
                                <InputAdornment position="start">
                                    <KeyboardIcon />
                                </InputAdornment>
                                ),
                            }}
                            sx={{ marginRight: '1em', marginBlock: '10px' }}
                            id="roomId"
                            label="Введите идентификатор комнаты"
                            value={roomId}
                            size="small"
                            onChange={(e) => handleRoomIdChange(e)}
                            onKeyUp={(event) => {
                                if (event.key== 'Enter') {
                                    handleJoinRoom();
                                }
                                    
                            }}
                        />
                        <SecondaryButton sx={{ minWidth: '0' }} disabled={!roomId.replace(/\s/g, '').length} onClick={handleJoinRoom} variant="text">
                            Join
                        </SecondaryButton>
                </Grid>
                <Grid item xs={12} md={6} sx={{ display: 'flex', justifyContent: 'center' }}>
                    {collabSvg()}
                </Grid>
            </Grid>            
            
        </>
    );
}