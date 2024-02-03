package entities.camera;

import org.lwjgl.util.vector.Vector3f;

public abstract class Camera {
    protected Vector3f position = new Vector3f(0, 0, 0);
    protected float pitch = 22.5f;
    protected float yaw = -225;
    protected float roll = 0;

    public abstract void move();

    public Vector3f getPosition() {
        return position;
    }

    public void setPosition(Vector3f position) {
        this.position = position;
    }

    public float getPitch() {
        return pitch;
    }

    public void setPitch(float pitch) {
        this.pitch = pitch;
    }

    public float getYaw() {
        return yaw;
    }

    public void setYaw(float yaw) {
        this.yaw = yaw;
    }

    public float getRoll() {
        return roll;
    }

    public void setRoll(float roll) {
        this.roll = roll;
    }
}
