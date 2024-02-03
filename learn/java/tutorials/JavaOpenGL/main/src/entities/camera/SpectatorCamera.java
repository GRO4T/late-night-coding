package entities.camera;

import org.lwjgl.input.Keyboard;
import org.lwjgl.input.Mouse;
import org.lwjgl.util.vector.Vector3f;

public class SpectatorCamera extends Camera{
    private Vector3f position = new Vector3f(0, 20, 0);
    private float pitch = 22.5f;
    private float yaw = -225;
    private float roll = 0;

    private final float cameraSpeed = 0.3f;

    @Override
    public void move(){
        if (Keyboard.isKeyDown(Keyboard.KEY_W)){
            position.x += Math.cos(Math.toRadians(pitch)) * Math.sin(Math.toRadians(yaw)) * cameraSpeed;
            position.y -= Math.sin(Math.toRadians(pitch)) * cameraSpeed;
            position.z -= Math.cos(Math.toRadians(pitch)) * Math.cos(Math.toRadians(yaw)) * cameraSpeed;
        }
        if (Keyboard.isKeyDown(Keyboard.KEY_S)){
            position.x -= Math.cos(Math.toRadians(pitch)) * Math.sin(Math.toRadians(yaw)) * cameraSpeed;
            position.y += Math.sin(Math.toRadians(pitch)) * cameraSpeed;
            position.z += Math.cos(Math.toRadians(pitch)) * Math.cos(Math.toRadians(yaw)) * cameraSpeed;
        }
        if (Keyboard.isKeyDown(Keyboard.KEY_D)){
            position.x += Math.cos(Math.toRadians(yaw)) * cameraSpeed;
            position.z += Math.sin(Math.toRadians(yaw)) * cameraSpeed;
        }
        if (Keyboard.isKeyDown(Keyboard.KEY_A)){
            position.x -= Math.cos(Math.toRadians(yaw)) * cameraSpeed;
            position.z -= Math.sin(Math.toRadians(yaw)) * cameraSpeed;
        }
        pitch -= Mouse.getDY() * cameraSpeed;
        yaw += Mouse.getDX() * cameraSpeed;

        //System.out.println("yaw: " + yaw);
        //System.out.println("pitch: " + pitch);
        //System.out.println(position);
    }
}
