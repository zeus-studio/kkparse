import React from 'react';
import styled from 'styled-components';
import { Input, Button, Typography, message, Row, Col, Space } from 'antd';
import axios, { AxiosProgressEvent } from 'axios';
import ReactPlayer from 'react-player';

import { Parse, Download } from "../wailsjs/go/main/App";
import type { parser } from '../wailsjs/go/models';

const { TextArea } = Input;
const { Title, Paragraph } = Typography;

const AppWrapper = styled.div`
    display: flex;
    flex-flow: column nowrap;
    align-items: center;
    .content {
        width: 100%;
    }
    .button {
        margin: 20px;
    }
    .player {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        video {
            max-height: 400px;
        }
    }
`;

function App() {
    const [urlVal, setUrlVal] = React.useState<string>('');
    const [loading, setLoading] = React.useState<boolean>(false);
    const [videoUrl, setVideoUrl] = React.useState<string>('');

    const handleParse = React.useCallback(async () => {
        if (!urlVal) {
            message.info('请先输入分享链接')
            return
        }
        try {
            setLoading(true)
            const res: parser.VideoParseInfo = await Parse(urlVal);
            setVideoUrl(res.video_url);
            setLoading(false)
        } catch (error) {
            message.warning('解析失败')
            setLoading(false)
        }
    }, [urlVal]);

    const [downloading, setDownloading] = React.useState<boolean>(false);
    const [downloadProgress, setDownloadProgress] = React.useState<string>('');
    const handleDownloadVideo = React.useCallback(async () => {
        setDownloading(true);
        await Download(videoUrl);
        setDownloading(false)
    }, []);

    return (
        <AppWrapper>
            <Typography>
                <Title>KK Parse</Title>
                <Paragraph>
                    支持解析抖音、快手、小红书、皮皮虾、火山视频、微博、微视、绿洲、最右、西瓜视频、梨视频、皮皮搞笑、虎牙、AcFun、逗拍、美拍、全民K歌、六间房、新片场、好看视频的视频资源。在视频软件的视频播放页面，点击分享按钮，复制分享链接。将链接输入到下面的
                    输入框里面，点击“解析”按钮即可尝试解析（可能会有解析失败的情况）
                </Paragraph>
            </Typography>
            <Row gutter={20} className="content">
                <Col span={12}>
                    <TextArea
                        className="textarea"
                        value={urlVal}
                        onChange={(e) => setUrlVal(e.target.value)}
                        autoSize={{ minRows: 6, maxRows: 12 }}
                        placeholder="在这里输入分享链接"
                    />
                    <Space wrap>
                        <Button className="button" type="primary" loading={loading} onClick={handleParse}>解析</Button>
                        <Button loading={downloading} disabled={!videoUrl} onClick={handleDownloadVideo}>下载{downloadProgress ? `（${downloadProgress}%）` : ''}</Button>
                    </Space>
                </Col>
                <Col span={12}>
                    <div className="player">
                        <ReactPlayer
                            url={videoUrl}
                            controls
                            width="auto"
                            height="100%"
                        />
                    </div>
                </Col>
            </Row>
        </AppWrapper>
    )
}

export default App
