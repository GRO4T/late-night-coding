package entities.camera;

import entities.Player;
import org.lwjgl.input.Mouse;
import toolbox.Maths;

public class ThirdPersonCamera extends Camera{
    private static float MIN_PITCH = 5;
    private static float MAX_PITCH = 90;
    private static float MIN_ZOOM = 20;
    private static float MAX_ZOOM = 400;

    private float distanceFromPlayer = 50;
    private float angleAroundPlayer = 0;

    private Player player;

    private float offsetY = 6;

    public ThirdPersonCamera(Player player) {
        this.player = player;
    }

    @Override
    public void move() {
        calculateZoom();
        calculatePitch();
        calculateAngleAroundPlayer();

        float horizontalDistance = calculateHorizontalDistance();
        float verticalDistance = calculateVerticalDistance();

        calculateCameraPosition(horizontalDistance, verticalDistance);

        //System.out.println("pitch: " + pitch);
        //System.out.println("zoom: " + distanceFromPlayer);
    }

    private void calculateZoom(){
        float zoomLevel = Mouse.getDWheel() * 0.1f;
        distanceFromPlayer -= zoomLevel;
        distanceFromPlayer = Maths.clamp(distanceFromPlayer, MIN_ZOOM, MAX_ZOOM);
    }

    private void calculatePitch(){
        // check if right mouse button is pressed
        if (Mouse.isButtonDown(1)){
            float pitchChange = Mouse.getDY() * 0.1f;
            pitch -= pitchChange;
            pitch = Maths.clamp(pitch, MIN_PITCH, MAX_PITCH);
        }
    }

    private void calculateAngleAroundPlayer(){
        // check if left mouse button is pressed
        if (Mouse.isButtonDown(0)){
            float angleChange = Mouse.getDX() * 0.3f;
            angleAroundPlayer -= angleChange;
        }
    }

    private float calculateHorizontalDistance(){
        return (float) (distanceFromPlayer * Math.cos(Math.toRadians(pitch)));
    }

    private float calculateVerticalDistance(){
        return (float) (distanceFromPlayer * Math.sin(Math.toRadians(pitch)));
    }

    private void calculateCameraPosition(float horizontalDistance, float verticalDistance){
        float theta = player.getRotY() + angleAroundPlayer;
        float offsetX = (float) (horizontalDistance * Math.sin(Math.toRadians(theta)));
        float offsetZ = (float) (horizontalDistance * Math.cos(Math.toRadians(theta)));
        position.x = player.getPosition().x - offsetX;
        position.y = player.getPosition().y + verticalDistance + offsetY;
        position.z = player.getPosition().z - offsetZ;
        yaw = 180.0f - theta;
    }

    public void setOffsetY(float offsetY) {
        this.offsetY = offsetY;
    }
}
